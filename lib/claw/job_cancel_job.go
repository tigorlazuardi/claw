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

// CancelJob cancels a running job
func (s *Claw) CancelJob(ctx context.Context, req *clawv1.CancelJobRequest) (*clawv1.CancelJobResponse, error) {
	nowMillis := types.UnixMilliNow()

	stmt := Jobs.UPDATE(Jobs.Status, Jobs.FinishedAt).
		MODEL(model.Jobs{
			Status:     clawv1.JobStatus_JOB_STATUS_CANCELLED.String(),
			FinishedAt: &nowMillis,
		}).
		WHERE(Jobs.ID.EQ(Int64(req.Id))).
		RETURNING(Jobs.AllColumns)

	var jobRow model.Jobs
	err := stmt.QueryContext(ctx, s.db, &jobRow)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel job: %w", err)
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

	return &clawv1.CancelJobResponse{
		Job: job,
	}, nil
}