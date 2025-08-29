package claw

import (
	"context"
	"fmt"
	"time"

	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	"github.com/tigorlazuardi/claw/lib/claw/types"
)

// StartScheduler begins the scheduler event loop
func (c *Claw) StartScheduler(ctx context.Context) {
	c.schedulerMutex.Lock()
	if c.schedulerRunning {
		c.schedulerMutex.Unlock()
		return
	}
	
	// Initialize channels if not already done
	if c.jobQueue == nil {
		c.jobQueue = make(chan *model.Jobs, 100)
	}
	if c.downloadQueue == nil {
		c.downloadQueue = make(chan downloadTask, 200)
	}
	if c.schedulerStopCh == nil {
		c.schedulerStopCh = make(chan struct{})
	}
	if c.schedulerDoneCh == nil {
		c.schedulerDoneCh = make(chan struct{})
	}
	
	c.schedulerRunning = true
	c.schedulerMutex.Unlock()
	
	defer close(c.schedulerDoneCh)
	defer func() {
		c.schedulerMutex.Lock()
		c.schedulerRunning = false
		c.schedulerMutex.Unlock()
	}()
	
	if c.Logger != nil {
		c.Logger.Info("Starting scheduler", "poll_interval", c.schedulerConfig.PollInterval, "max_workers", c.schedulerConfig.MaxWorkers, "download_workers", c.schedulerConfig.DownloadWorkers)
	}

	// Start job workers
	for i := 0; i < c.schedulerConfig.MaxWorkers; i++ {
		go c.jobWorker(ctx, i)
	}

	// Start download workers
	for i := 0; i < c.schedulerConfig.DownloadWorkers; i++ {
		go c.downloadWorker(ctx, i)
	}

	// Start job polling
	ticker := time.NewTicker(c.schedulerConfig.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if c.Logger != nil {
				c.Logger.Info("Scheduler context cancelled, stopping")
			}
			return
		case <-c.schedulerStopCh:
			if c.Logger != nil {
				c.Logger.Info("Scheduler stop signal received, stopping")
			}
			return
		case <-ticker.C:
			if err := c.pollJobs(ctx); err != nil && c.Logger != nil {
				c.Logger.Error("Failed to poll jobs", "error", err)
			}
		}
	}
}

// StopScheduler gracefully stops the scheduler
func (c *Claw) StopScheduler() {
	c.schedulerMutex.RLock()
	if !c.schedulerRunning || c.schedulerStopCh == nil {
		c.schedulerMutex.RUnlock()
		return
	}
	c.schedulerMutex.RUnlock()
	
	close(c.schedulerStopCh)
}

// WaitScheduler waits for the scheduler to finish
func (c *Claw) WaitScheduler() {
	c.schedulerMutex.RLock()
	doneCh := c.schedulerDoneCh
	c.schedulerMutex.RUnlock()
	
	if doneCh != nil {
		<-doneCh
	}
}

// pollJobs polls the database for pending jobs and queues them
func (c *Claw) pollJobs(ctx context.Context) error {
	// Use existing ListJobs API to get pending jobs
	resp, err := c.ListJobs(ctx, &clawv1.ListJobsRequest{
		Status: clawv1.JobStatus_JOB_STATUS_PENDING.Enum(),
		Sorts: []*clawv1.ListJobsRequest_Sort{{
			Field: clawv1.JobSortField_JOB_SORT_FIELD_CREATED_AT,
			Desc:  false,
		}},
		Pagination: &clawv1.Pagination{
			Size: Ptr(uint32(100)),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to query pending jobs: %w", err)
	}

	c.queuedJobsMutex.Lock()
	defer c.queuedJobsMutex.Unlock()

	for _, job := range resp.Jobs {
		// Skip if already queued
		if c.queuedJobs[job.Id] {
			continue
		}

		// Convert proto job back to model for internal processing
		jobModel := &model.Jobs{
			ID:       &job.Id,
			SourceID: job.SourceId,
			Status:   job.Status.String(),
			CreatedAt: types.NewUnixMilli(job.CreatedAt.AsTime()),
		}
		if job.ScheduleId != nil {
			jobModel.ScheduleID = *job.ScheduleId
		}
		if job.RunAt != nil {
			runAt := types.NewUnixMilli(job.RunAt.AsTime())
			jobModel.RunAt = &runAt
		}
		if job.FinishedAt != nil {
			finishedAt := types.NewUnixMilli(job.FinishedAt.AsTime())
			jobModel.FinishedAt = &finishedAt
		}
		if job.Error != nil {
			jobModel.Error = job.Error
		}

		// Mark as queued and send to job queue
		c.queuedJobs[job.Id] = true
		
		select {
		case c.jobQueue <- jobModel:
			if c.Logger != nil {
				c.Logger.Debug("Queued job", "job_id", job.Id, "source_id", job.SourceId)
			}
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Job queue is full, will try again on next poll
			delete(c.queuedJobs, job.Id)
			if c.Logger != nil {
				c.Logger.Warn("Job queue full, skipping job", "job_id", job.Id)
			}
		}
	}

	return nil
}

// removeFromQueue removes a job from the queued jobs map
func (c *Claw) removeFromQueue(jobID int64) {
	c.queuedJobsMutex.Lock()
	defer c.queuedJobsMutex.Unlock()
	delete(c.queuedJobs, jobID)
}