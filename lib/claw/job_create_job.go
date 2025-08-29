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

// CreateJob creates a new job
func (s *Claw) CreateJob(ctx context.Context, req *clawv1.CreateJobRequest) (*clawv1.CreateJobResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	nowMillis := types.UnixMilliNow()

	// Insert job
	columns := ColumnList{
		Jobs.SourceID,
		Jobs.Status,
		Jobs.CreatedAt,
	}
	jobModel := model.Jobs{
		SourceID:  req.SourceId,
		Status:    req.Status.String(),
		CreatedAt: nowMillis,
	}
	
	// Only include schedule_id if provided
	if req.ScheduleId != nil {
		columns = append(columns, Jobs.ScheduleID)
		jobModel.ScheduleID = *req.ScheduleId
	}

	jobStmt := Jobs.INSERT(columns).MODEL(jobModel).RETURNING(Jobs.AllColumns)

	var jobRow model.Jobs
	err = jobStmt.QueryContext(ctx, tx, &jobRow)
	if err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	// Create job images if provided
	if len(req.JobImages) > 0 {
		var jobImages []model.JobImages
		for _, jobImage := range req.JobImages {
			jobImages = append(jobImages, model.JobImages{
				JobID:     *jobRow.ID,
				ImageID:   jobImage.ImageId,
				DeviceID:  jobImage.DeviceId,
				Action:    jobImage.Action.String(),
				CreatedAt: nowMillis,
			})
		}

		jobImageStmt := JobImages.
			INSERT(JobImages.JobID, JobImages.ImageID, JobImages.DeviceID, JobImages.Action, JobImages.CreatedAt).
			MODELS(jobImages)

		_, err = jobImageStmt.ExecContext(ctx, tx)
		if err != nil {
			return nil, fmt.Errorf("failed to create job images: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Convert to protobuf
	job := &clawv1.Job{
		Id:        *jobRow.ID,
		SourceId:  jobRow.SourceID,
		Status:    clawv1.JobStatus(clawv1.JobStatus_value[jobRow.Status]),
		CreatedAt: jobRow.CreatedAt.ToProto(),
	}

	if jobRow.ScheduleID != 0 {
		job.ScheduleId = &jobRow.ScheduleID
	}
	if jobRow.RunAt != nil {
		job.RunAt = jobRow.RunAt.ToProto()
	}
	if jobRow.FinishedAt != nil {
		job.FinishedAt = jobRow.FinishedAt.ToProto()
	}
	if jobRow.Error != nil {
		job.Error = jobRow.Error
	}

	// Add job images
	for _, jobImageReq := range req.JobImages {
		job.JobImages = append(job.JobImages, &clawv1.JobImage{
			Id:        0, // Will be set by database
			JobId:     *jobRow.ID,
			ImageId:   jobImageReq.ImageId,
			DeviceId:  jobImageReq.DeviceId,
			Action:    jobImageReq.Action,
			CreatedAt: nowMillis.ToProto(),
		})
	}

	return &clawv1.CreateJobResponse{
		Job: job,
	}, nil
}