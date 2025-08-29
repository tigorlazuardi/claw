package claw

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/claw/source"
)

// jobWorker processes jobs from the job queue
func (c *Claw) jobWorker(ctx context.Context, workerID int) {
	var logger *slog.Logger
	if c.logger != nil {
		logger = c.logger.With("worker_id", workerID, "worker_type", "job")
		logger.Debug("Starting job worker")
		defer logger.Debug("Job worker stopped")
	}

	for {
		select {
		case <-ctx.Done():
			return
		case job := <-c.jobQueue:
			if job.ID == nil {
				continue
			}

			if logger != nil {
				logger = logger.With("job_id", *job.ID)
				logger.Debug("Processing job")
			}

			if err := c.processJob(ctx, job); err != nil {
				if logger != nil {
					logger.Error("Failed to process job", "error", err)
				}
				c.UpdateJob(ctx, &clawv1.UpdateJobRequest{
					Id:     *job.ID,
					Status: clawv1.JobStatus_JOB_STATUS_FAILED.Enum(),
					Error:  Ptr(err.Error()),
				})
			}

			// Remove from queued jobs regardless of success/failure
			c.removeFromQueue(*job.ID)
		}
	}
}

// processJob processes a single job
func (c *Claw) processJob(ctx context.Context, job *model.Jobs) error {
	// Update job status to running using existing API
	_, err := c.UpdateJob(ctx, &clawv1.UpdateJobRequest{
		Id:     *job.ID,
		Status: clawv1.JobStatus_JOB_STATUS_RUNNING.Enum(),
	})
	if err != nil {
		return fmt.Errorf("failed to update job status to running: %w", err)
	}

	// Get source details using existing API
	sourceResp, err := c.GetSource(ctx, &clawv1.GetSourceRequest{
		Id: job.SourceID,
	})
	if err != nil {
		return fmt.Errorf("failed to get source: %w", err)
	}

	src := sourceResp.Source

	// Get the source implementation
	backend, exists := c.sources[src.Kind]
	if !exists {
		return fmt.Errorf("unknown source kind: %s", src.Kind)
	}

	// Prepare source request
	countback := 25 // Default value
	if src.Countback > 0 {
		countback = int(src.Countback)
	}

	sourceReq := source.Request{
		Parameter: src.Parameter,
		Countback: countback,
	}

	if c.logger != nil {
		c.logger.Debug("Executing source backend", "source_kind", src.Kind, "parameter", src.Parameter, "countback", countback)
	}

	// Execute the source
	resp, err := backend.Run(ctx, sourceReq)
	if err != nil {
		return fmt.Errorf("source execution failed: %w", err)
	}

	if c.logger != nil {
		c.logger.Info("Source execution completed", "images_found", len(resp.Images))
	}

	if len(resp.Images) == 0 {
		// No images found, mark job as completed using existing API
		_, err := c.UpdateJob(ctx, &clawv1.UpdateJobRequest{
			Id:     *job.ID,
			Status: clawv1.JobStatus_JOB_STATUS_COMPLETED.Enum(),
		})
		if err != nil && c.logger != nil {
			c.logger.Error("Failed to update job status to completed", "error", err)
		}
		return nil
	}

	// Get devices for filtering
	devices, err := c.getDevicesForJob(ctx, *job.ID)
	if err != nil {
		return fmt.Errorf("failed to get devices for job: %w", err)
	}

	// Process images
	for _, img := range resp.Images {
		// Filter devices that match this image
		matchedDevices := c.filterDevicesForImage(img, devices)
		if len(matchedDevices) == 0 {
			if c.logger != nil {
				c.logger.Debug("No devices match image, skipping", "download_url", img.DownloadURL)
			}
			continue
		}

		// Queue for download
		task := downloadTask{
			jobID:      *job.ID,
			sourceID:   job.SourceID,
			image:      img,
			devices:    matchedDevices,
			sourceName: src.Kind,
		}

		select {
		case c.downloadQueue <- task:
			if c.logger != nil {
				c.logger.Debug("Queued image for download", "download_url", img.DownloadURL, "matched_devices", len(matchedDevices))
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// Mark job as completed using existing API
	_, err = c.UpdateJob(ctx, &clawv1.UpdateJobRequest{
		Id:     *job.ID,
		Status: clawv1.JobStatus_JOB_STATUS_COMPLETED.Enum(),
	})
	if err != nil && c.logger != nil {
		c.logger.Error("Failed to update job status to completed", "error", err)
	}

	return nil
}

