package claw

import (
	"context"
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
)

// GetJob retrieves a job by ID
func (s *Claw) GetJob(ctx context.Context, req *clawv1.GetJobRequest) (*clawv1.GetJobResponse, error) {
	cols := ProjectionList{Jobs.AllColumns}
	var from ReadableTable = Jobs
	if req.GetIncludeJobImages() {
		cols = append(cols, JobImages.AllColumns)
		from = Jobs.LEFT_JOIN(JobImages, JobImages.JobID.EQ(Jobs.ID))
	}
	query := SELECT(cols).
		FROM(from).
		WHERE(Jobs.ID.EQ(Int64(req.Id)))

	var out struct {
		model.Jobs
		JobImages []model.JobImages
	}
	err := query.QueryContext(ctx, s.db, &out)
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	// Convert to protobuf
	job := &clawv1.Job{
		Id:        *out.ID,
		SourceId:  out.SourceID,
		Status:    clawv1.JobStatus(clawv1.JobStatus_value[out.Status]),
		CreatedAt: out.CreatedAt.ToProto(),
	}

	if out.ScheduleID != 0 {
		job.ScheduleId = &out.ScheduleID
	}
	if out.RunAt != nil {
		job.RunAt = out.RunAt.ToProto()
	}
	if out.FinishedAt != nil {
		job.FinishedAt = out.FinishedAt.ToProto()
	}
	if out.Error != nil {
		job.Error = out.Error
	}

	return &clawv1.GetJobResponse{
		Job: job,
	}, nil
}

