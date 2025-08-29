package claw

import (
	"context"
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
)

// ListJobs lists jobs with optional filtering and pagination
func (s *Claw) ListJobs(ctx context.Context, req *clawv1.ListJobsRequest) (*clawv1.ListJobsResponse, error) {
	var from ReadableTable = Jobs
	
	// Set up FROM clause with optional join for device filtering
	if req.DeviceId != nil {
		from = Jobs.LEFT_JOIN(JobImages, JobImages.JobID.EQ(Jobs.ID))
	}
	
	query := SELECT(Jobs.AllColumns).FROM(from)

	// Apply filters
	var conditions []BoolExpression
	if req.SourceId != nil {
		conditions = append(conditions, Jobs.SourceID.EQ(Int64(*req.SourceId)))
	}
	if req.ScheduleId != nil {
		conditions = append(conditions, Jobs.ScheduleID.EQ(Int64(*req.ScheduleId)))
	}
	if req.Status != nil {
		conditions = append(conditions, Jobs.Status.EQ(String(req.Status.String())))
	}
	if req.DeviceId != nil {
		conditions = append(conditions, JobImages.DeviceID.EQ(Int64(*req.DeviceId)))
	}

	if len(conditions) > 0 {
		query = query.WHERE(AND(conditions...))
	}

	// Apply sorting
	if len(req.Sorts) > 0 {
		var orderBy []OrderByClause
		for _, sort := range req.Sorts {
			var column Column
			switch sort.Field {
			case clawv1.JobSortField_JOB_SORT_FIELD_ID:
				column = Jobs.ID
			case clawv1.JobSortField_JOB_SORT_FIELD_SOURCE_ID:
				column = Jobs.SourceID
			case clawv1.JobSortField_JOB_SORT_FIELD_SCHEDULE_ID:
				column = Jobs.ScheduleID
			case clawv1.JobSortField_JOB_SORT_FIELD_CREATED_AT:
				column = Jobs.CreatedAt
			case clawv1.JobSortField_JOB_SORT_FIELD_RUN_AT:
				column = Jobs.RunAt
			case clawv1.JobSortField_JOB_SORT_FIELD_FINISHED_AT:
				column = Jobs.FinishedAt
			case clawv1.JobSortField_JOB_SORT_FIELD_STATUS:
				column = Jobs.Status
			default:
				column = Jobs.ID
			}

			if sort.Desc {
				orderBy = append(orderBy, column.DESC())
			} else {
				orderBy = append(orderBy, column.ASC())
			}
		}
		query = query.ORDER_BY(orderBy...)
	} else {
		// Default sorting by ID descending
		query = query.ORDER_BY(Jobs.ID.DESC())
	}

	// Apply pagination
	if req.Pagination != nil {
		if req.Pagination.Size != nil && *req.Pagination.Size > 0 {
			query = query.LIMIT(int64(*req.Pagination.Size))
		}
		// Token-based pagination would need additional logic here
	}

	var jobRows []model.Jobs
	err := query.QueryContext(ctx, s.db, &jobRows)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}

	// Convert to protobuf
	var jobs []*clawv1.Job
	for _, jobRow := range jobRows {
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

		// Load job images if requested
		if req.IncludeJobImages != nil && *req.IncludeJobImages {
			jobImagesQuery := SELECT(JobImages.AllColumns).
				FROM(JobImages).
				WHERE(JobImages.JobID.EQ(Int64(*jobRow.ID)))

			var jobImageRows []model.JobImages
			err = jobImagesQuery.QueryContext(ctx, s.db, &jobImageRows)
			if err != nil {
				return nil, fmt.Errorf("failed to get job images: %w", err)
			}

			for _, jobImageRow := range jobImageRows {
				job.JobImages = append(job.JobImages, &clawv1.JobImage{
					Id:        *jobImageRow.ID,
					JobId:     jobImageRow.JobID,
					ImageId:   jobImageRow.ImageID,
					DeviceId:  jobImageRow.DeviceID,
					Action:    clawv1.JobAction(clawv1.JobAction_value[jobImageRow.Action]),
					CreatedAt: jobImageRow.CreatedAt.ToProto(),
				})
			}
		}

		jobs = append(jobs, job)
	}

	// Create response pagination
	responsePagination := &clawv1.Pagination{}
	if req.Pagination != nil {
		responsePagination = req.Pagination
	}

	return &clawv1.ListJobsResponse{
		Jobs:       jobs,
		Pagination: responsePagination,
	}, nil
}