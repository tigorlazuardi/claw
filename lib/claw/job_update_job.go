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

// UpdateJob updates an existing job
func (s *Claw) UpdateJob(ctx context.Context, req *clawv1.UpdateJobRequest) (*clawv1.UpdateJobResponse, error) {
	nowMillis := types.UnixMilliNow()

	columns := ColumnList{}
	updateModel := model.Jobs{}

	if req.Status != nil {
		columns = append(columns, Jobs.Status)
		updateModel.Status = req.Status.String()

		// Set run_at when status changes to RUNNING
		if *req.Status == clawv1.JobStatus_JOB_STATUS_RUNNING {
			columns = append(columns, Jobs.RunAt)
			updateModel.RunAt = &nowMillis
		}

		// Set finished_at when status changes to COMPLETED, FAILED, or CANCELLED
		if *req.Status == clawv1.JobStatus_JOB_STATUS_COMPLETED ||
			*req.Status == clawv1.JobStatus_JOB_STATUS_FAILED ||
			*req.Status == clawv1.JobStatus_JOB_STATUS_CANCELLED {
			columns = append(columns, Jobs.FinishedAt)
			updateModel.FinishedAt = &nowMillis
		}
	}

	if req.Error != nil {
		columns = append(columns, Jobs.Error)
		updateModel.Error = req.Error
	}

	stmt := Jobs.UPDATE(columns).
		MODEL(updateModel).
		WHERE(Jobs.ID.EQ(Int64(req.Id))).
		RETURNING(Jobs.AllColumns)

	var jobRow model.Jobs
	err := stmt.QueryContext(ctx, s.db, &jobRow)
	if err != nil {
		return nil, fmt.Errorf("failed to update job: %w", err)
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

	return &clawv1.UpdateJobResponse{
		Job: job,
	}, nil
}