package claw

import (
	"context"
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/claw/types"
)

// RetryJob creates a new job based on a failed job
func (s *Claw) RetryJob(ctx context.Context, req *clawv1.RetryJobRequest) (*clawv1.RetryJobResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get the original job
	originalJobQuery := SELECT(Jobs.AllColumns).
		FROM(Jobs).
		WHERE(Jobs.ID.EQ(Int64(req.Id)))

	var originalJob model.Jobs
	err = originalJobQuery.QueryContext(ctx, tx, &originalJob)
	if err != nil {
		return nil, fmt.Errorf("failed to get original job: %w", err)
	}

	// Get the original job images
	jobImagesQuery := SELECT(JobImages.AllColumns).
		FROM(JobImages).
		WHERE(JobImages.JobID.EQ(Int64(req.Id)))

	var originalJobImages []model.JobImages
	err = jobImagesQuery.QueryContext(ctx, tx, &originalJobImages)
	if err != nil {
		return nil, fmt.Errorf("failed to get original job images: %w", err)
	}

	nowMillis := types.UnixMilliNow()

	// Create new job
	newJobStmt := Jobs.INSERT(
		Jobs.SourceID,
		Jobs.ScheduleID,
		Jobs.Status,
		Jobs.CreatedAt,
	).MODEL(model.Jobs{
		SourceID:   originalJob.SourceID,
		ScheduleID: originalJob.ScheduleID,
		Status:     clawv1.JobStatus_JOB_STATUS_PENDING.String(),
		CreatedAt:  nowMillis,
	}).RETURNING(Jobs.AllColumns)

	var newJobRow model.Jobs
	err = newJobStmt.QueryContext(ctx, tx, &newJobRow)
	if err != nil {
		return nil, fmt.Errorf("failed to create retry job: %w", err)
	}

	// Create new job images
	if len(originalJobImages) > 0 {
		var newJobImages []model.JobImages
		for _, originalJobImage := range originalJobImages {
			newJobImages = append(newJobImages, model.JobImages{
				JobID:     *newJobRow.ID,
				ImageID:   originalJobImage.ImageID,
				DeviceID:  originalJobImage.DeviceID,
				Action:    originalJobImage.Action,
				CreatedAt: nowMillis,
			})
		}

		jobImageStmt := JobImages.
			INSERT(JobImages.JobID, JobImages.ImageID, JobImages.DeviceID, JobImages.Action, JobImages.CreatedAt).
			MODELS(newJobImages)

		_, err = jobImageStmt.ExecContext(ctx, tx)
		if err != nil {
			return nil, fmt.Errorf("failed to create retry job images: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Convert to protobuf
	job := &clawv1.Job{
		Id:        *newJobRow.ID,
		SourceId:  newJobRow.SourceID,
		Status:    clawv1.JobStatus(clawv1.JobStatus_value[newJobRow.Status]),
		CreatedAt: newJobRow.CreatedAt.ToProto(),
	}

	if newJobRow.ScheduleID != 0 {
		job.ScheduleId = &newJobRow.ScheduleID
	}
	if newJobRow.RunAt != nil {
		job.RunAt = newJobRow.RunAt.ToProto()
	}
	if newJobRow.FinishedAt != nil {
		job.FinishedAt = newJobRow.FinishedAt.ToProto()
	}
	if newJobRow.Error != nil {
		job.Error = newJobRow.Error
	}

	// Add job images to response
	for _, originalJobImage := range originalJobImages {
		job.JobImages = append(job.JobImages, &clawv1.JobImage{
			Id:        0, // Will be set by database
			JobId:     *newJobRow.ID,
			ImageId:   originalJobImage.ImageID,
			DeviceId:  originalJobImage.DeviceID,
			Action:    clawv1.JobAction(clawv1.JobAction_value[originalJobImage.Action]),
			CreatedAt: nowMillis.ToProto(),
		})
	}

	return &clawv1.RetryJobResponse{
		Job: job,
	}, nil
}