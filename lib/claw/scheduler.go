package claw

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/teivah/broadcast"
	"github.com/tigorlazuardi/claw/lib/claw/config"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/claw/source"
	"github.com/tigorlazuardi/claw/lib/claw/types"
	"github.com/tigorlazuardi/claw/lib/dblogger"
	"golang.org/x/sync/semaphore"
)

const leastCommonMultiple = 720720 // LCM of 1 to 16

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type scheduler struct {
	claw           *Claw
	config         *config.Config
	isRunning      atomic.Bool
	tracker        *tracker
	queue          chan model.Jobs
	imageSemaphore *semaphore.Weighted
	wg             *sync.WaitGroup
	reloadSignal   *broadcast.Relay[struct{}]
	logger         *slog.Logger
	backends       map[string]source.Source
	httpclient     Doer
}

type imageQueue struct {
	image   source.Image
	devices []model.Devices
}

func (scheduler *scheduler) start(baseContext context.Context) {
	if scheduler.isRunning.Load() {
		return
	}
	scheduler.isRunning.Store(true)
	defer scheduler.isRunning.Store(false)
	go scheduler.startPolling(baseContext)
	go scheduler.consumeJobQueue(baseContext)
	scheduler.logger.Info("scheduler started")
	<-baseContext.Done()
	scheduler.logger.Info("shutting down scheduler, waiting for running jobs to complete")
	ctx, cancel := context.WithTimeout(context.Background(), scheduler.config.Scheduler.ExitTimeout)
	defer cancel()
	wait := make(chan struct{}, 1)
	go func() {
		scheduler.wg.Wait()
		wait <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		scheduler.logger.Warn("scheduler shutdown timed out, some jobs may be interrupted")
	case <-wait:
		scheduler.logger.Info("scheduler shutdown complete")
	}
}

func (scheduler *scheduler) startPolling(ctx context.Context) {
	if jobs, _ := scheduler.getJobs(ctx); len(jobs) > 0 {
		scheduler.enqueueJobs(jobs)
	}
	ticker := time.NewTicker(scheduler.config.Scheduler.PollInterval)
	defer ticker.Stop()
	reload := scheduler.reloadSignal.Listener(1)
	defer reload.Close()

	for {
		select {
		case <-ctx.Done():
			scheduler.logger.DebugContext(ctx, "scheduler poller stopped")
			return
		case <-reload.Ch():
			scheduler.logger.InfoContext(ctx, "reloading scheduler poll interval", "new_interval", scheduler.config.Scheduler.PollInterval)
			ticker.Reset(scheduler.config.Scheduler.PollInterval)
		case <-ticker.C:
			jobs, err := scheduler.getJobs(ctx)
			if err != nil {
				scheduler.logger.ErrorContext(ctx, "failed to get jobs", "error", err)
				continue
			}
			scheduler.enqueueJobs(jobs)
		}
	}
}

func (scheduler *scheduler) getJobs(ctx context.Context) ([]model.Jobs, error) {
	var jobs []model.Jobs
	cond := Jobs.FinishedAt.IS_NULL()
	if runningIds := scheduler.tracker.List(); len(runningIds) > 0 {
		expr := make([]Expression, len(runningIds))
		for i, id := range runningIds {
			expr[i] = Int(id)
		}
		cond = AND(cond, Jobs.ID.NOT_IN(expr...))
	}
	ctx = dblogger.ContextWithSkipLog(ctx)
	err := SELECT(Jobs.AllColumns).
		FROM(Jobs).
		WHERE(cond).
		ORDER_BY(Jobs.CreatedAt.ASC()).
		QueryContext(ctx, scheduler.claw.db, &jobs)
	if err != nil {
		return nil, fmt.Errorf("failed to query jobs: %w", err)
	}
	return jobs, nil
}

func (scheduler *scheduler) enqueueJobs(jobs []model.Jobs) {
	for _, job := range jobs {
		scheduler.tracker.Add(*job.ID)
		scheduler.logger.Info("enqueuing job", "job_id", job.ID, "source_id", job.SourceID, "schedule_id", job.ScheduleID)
		scheduler.queue <- job
		scheduler.logger.Info("job run", "job_id", job.ID, "source_id", job.SourceID, "schedule_id", job.ScheduleID)
	}
}

func (scheduler *scheduler) consumeJobQueue(ctx context.Context) {
	sem := semaphore.NewWeighted(leastCommonMultiple)
	workers := min(scheduler.config.Scheduler.MaxWorkers, 16) // Max 16 workers
	weight := int64(leastCommonMultiple / workers)
	reload := scheduler.reloadSignal.Listener(1)
	for {
		select {
		case <-ctx.Done():
			scheduler.logger.DebugContext(ctx, "scheduler queue consumer stopped")
			return
		case <-reload.Ch():
			workers = min(scheduler.config.Scheduler.MaxWorkers, 16)
			weight = int64(leastCommonMultiple / workers)
			scheduler.logger.InfoContext(ctx, "reloading max workers", "new_max_workers", workers)
		case job := <-scheduler.queue:
			currentWeight := weight // copy current weight before moving to goroutine
			if err := sem.Acquire(ctx, weight); err != nil {
				// Context canceled
				scheduler.logger.DebugContext(ctx, "scheduler queue consumer stopped")
				return
			}
			scheduler.wg.Add(1)
			go func(job model.Jobs) {
				defer func() {
					sem.Release(currentWeight)
					scheduler.wg.Done()
					scheduler.tracker.Remove(*job.ID)
				}()
				scheduler.executeJob(ctx, *job.ID)
			}(job)
		}
	}
}

func (scheduler *scheduler) executeJob(ctx context.Context, job int64) {
	var (
		src model.Sources
		err error
	)
	err = SELECT(Sources.AllColumns).
		WHERE(Sources.ID.EQ(Int64(job))).
		QueryContext(ctx, scheduler.claw.db, &src)
	if err != nil {
		scheduler.logger.ErrorContext(ctx, "failed to get source for job", "job_id", job, "error", err)
		return
	}

	backend, ok := scheduler.backends[src.Name]
	if !ok {
		err = fmt.Errorf("no backend found for source name: %s", src.Name)
		scheduler.logger.ErrorContext(ctx, "failed to get backend for job", "job_id", job, "source_name", src.Name, "error", err)
		scheduler.updateJobStatus(ctx, job, clawv1.JobStatus_JOB_STATUS_FAILED, updateJobStatusAttributes{
			err:        err,
			finishedAt: Ptr(types.UnixMilliNow()),
		})
		return
	}
	scheduler.updateJobStatus(ctx, job, clawv1.JobStatus_JOB_STATUS_RUNNING, updateJobStatusAttributes{
		runAt: Ptr(types.UnixMilliNow()),
	})
	scheduler.logger.InfoContext(ctx, "starting job", "job_id", job, "source_id", src.ID, "source_name", src.Name)

	resp, err := backend.Run(ctx, source.Request{
		Parameter: src.Parameter,
		Countback: int(src.Countback),
	})
	if err != nil {
		scheduler.logger.ErrorContext(ctx, "job failed", "job_id", job, "error", err)
		scheduler.updateJobStatus(ctx, job, clawv1.JobStatus_JOB_STATUS_FAILED, updateJobStatusAttributes{
			err:        err,
			finishedAt: Ptr(types.UnixMilliNow()),
		})
		return
	}
	if len(resp.Images) == 0 {
		scheduler.logger.InfoContext(ctx, "job completed with no images", "job_id", job)
		scheduler.updateJobStatus(ctx, job, clawv1.JobStatus_JOB_STATUS_COMPLETED, updateJobStatusAttributes{
			finishedAt: Ptr(types.UnixMilliNow()),
		})
		return
	}
	wg := sync.WaitGroup{}
	completed := make([]imageQueue, len(resp.Images))
	for i, image := range resp.Images {
		devices, err := scheduler.findDevicesToAssign(ctx, image)
		if err != nil {
			scheduler.logger.ErrorContext(ctx, "failed to find devices to assign", "job_id", job, "error", err)
			scheduler.updateJobStatus(ctx, job, clawv1.JobStatus_JOB_STATUS_FAILED, updateJobStatusAttributes{
				err:        err,
				finishedAt: Ptr(types.UnixMilliNow()),
			})
			return
		}
		if len(devices) == 0 {
			scheduler.logger.InfoContext(ctx, "no devices found to assign image", "job_id", job, "image", image)
			continue
		}
		wg.Add(1)
		weight := leastCommonMultiple / min(int64(scheduler.config.Scheduler.DownloadWorkers), 16)
		if err := scheduler.imageSemaphore.Acquire(ctx, weight); err != nil {
			// context canceled
			return
		}
		go func(image source.Image, devices []model.Devices) {
			defer wg.Done()
			defer scheduler.imageSemaphore.Release(weight)
			if err := scheduler.processDownload(ctx, image, devices, src.Name); err != nil {
				scheduler.logger.ErrorContext(ctx, "failed to process image", "job_id", job, "image", image, "error", err)
				return
			}
			completed[i] = imageQueue{image: image, devices: devices}
		}(image, devices)
	}
	wg.Wait()
	scheduler.logger.InfoContext(ctx, "job completed", "job_id", job, "images_processed", len(resp.Images))
	scheduler.updateJobStatus(ctx, job, clawv1.JobStatus_JOB_STATUS_COMPLETED, updateJobStatusAttributes{
		finishedAt: Ptr(types.UnixMilliNow()),
		collectedImages: slices.DeleteFunc(completed, func(queue imageQueue) bool {
			return queue.image.DownloadURL == "" // filter out invalid data
		}),
	})
}

func (scheduler *scheduler) findDevicesToAssign(ctx context.Context, image source.Image) ([]model.Devices, error) {
	imageRatio := float64(image.Width) / float64(image.Height)
	cond := Devices.IsEnabled.EQ(Int(1)).
		AND(
			Float(imageRatio).BETWEEN(
				CAST(Devices.Width).AS_REAL().DIV(CAST(Devices.Width).AS_REAL()).SUB(Devices.AspectRatioDifference),
				CAST(Devices.Width).AS_REAL().DIV(CAST(Devices.Width).AS_REAL()).ADD(Devices.AspectRatioDifference),
			),
		).
		AND(
			Devices.ImageMinWidth.LT_EQ(Int(0)).OR(Devices.ImageMinWidth.LT_EQ(Int(image.Width))),
		).
		AND(
			Devices.ImageMaxWidth.LT_EQ(Int(0)).OR(Devices.ImageMaxWidth.GT_EQ(Int(image.Width))),
		).
		AND(
			Devices.ImageMinHeight.LT_EQ(Int(0)).OR(Devices.ImageMinWidth.LT_EQ(Int(image.Height))),
		).
		AND(
			Devices.ImageMaxHeight.LT_EQ(Int(0)).OR(Devices.ImageMaxWidth.GT_EQ(Int(image.Height))),
		).
		AND(
			Devices.ImageMinFileSize.LT_EQ(Int(0)).OR(Devices.ImageMinFileSize.LT_EQ(Int(image.Filesize))),
		).
		AND(
			Devices.ImageMaxFileSize.LT_EQ(Int(0)).OR(Devices.ImageMaxFileSize.GT_EQ(Int(image.Filesize))),
		)
	if image.NSFW {
		cond = AND(cond, Devices.NsfwMode.NOT_EQ(Int(2)))
	} else {
		cond = AND(cond, Devices.NsfwMode.NOT_EQ(Int(3)))
	}
	var devices []model.Devices
	err := SELECT(Devices.AllColumns).
		FROM(Devices).
		WHERE(cond).
		QueryContext(ctx, scheduler.claw.db, &devices)
	if err != nil {
		return nil, fmt.Errorf("failed to query devices to assign: %w", err)
	}
	return devices, nil
}

type updateJobStatusAttributes struct {
	err             error
	runAt           *types.UnixMilli
	finishedAt      *types.UnixMilli
	collectedImages []imageQueue
	failedImages    []imageQueue
}

func (scheduler *scheduler) updateJobStatus(ctx context.Context, job int64, status clawv1.JobStatus, attr updateJobStatusAttributes) {
	// ContextCancelled error should only happens when the job is cancelled not by user.
	//
	// Graceful exits must not update the job status to failed.
	if errors.Is(attr.err, context.Canceled) || errors.Is(attr.err, context.DeadlineExceeded) {
		return
	}
	value := model.Jobs{}
	col := ColumnList{}
	if attr.err != nil {
		value.Status = clawv1.JobStatus_JOB_STATUS_FAILED.String()
		value.Error = Ptr(attr.err.Error())
		col = append(col, Jobs.Status, Jobs.Error)
	} else {
		value.Status = status.String()
		col = append(col, Jobs.Status)
	}
	if attr.finishedAt != nil {
		value.FinishedAt = attr.finishedAt
		col = append(col, Jobs.FinishedAt)
	}
	if attr.runAt != nil {
		value.RunAt = attr.runAt
		col = append(col, Jobs.RunAt)
	}
	_, err := Jobs.
		UPDATE(col).
		MODEL(value).
		WHERE(Jobs.ID.EQ(Int64(job))).
		ExecContext(ctx, scheduler.claw.db)
	if err != nil {
		scheduler.logger.ErrorContext(ctx, "failed to update job status", "job_id", job, "error", err, "status", status.String())
	}
}
