package claw

import (
	"context"
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
)

// DeleteJobs deletes jobs by IDs
func (s *Claw) DeleteJobs(ctx context.Context, req *clawv1.DeleteJobsRequest) (*clawv1.DeleteJobsResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	ids := make([]Expression, len(req.Ids))
	for i, id := range req.Ids {
		ids[i] = Int64(id)
	}

	// Delete job images first
	jobImagesStmt := JobImages.DELETE().WHERE(JobImages.JobID.IN(ids...))
	_, err = jobImagesStmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to delete job images: %w", err)
	}

	// Delete jobs
	jobsStmt := Jobs.DELETE().WHERE(Jobs.ID.IN(ids...))
	result, err := jobsStmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to delete jobs: %w", err)
	}

	deletedCount, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get deleted count: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &clawv1.DeleteJobsResponse{
		DeletedCount: int32(deletedCount),
	}, nil
}