package claw

import (
	"context"
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
)

// DeleteSource deletes a source by ID
func (s *Claw) DeleteSource(ctx context.Context, req *clawv1.DeleteSourceRequest) (*clawv1.DeleteSourceResponse, error) {
	in := make([]Expression, 0, len(req.Ids))
	for _, id := range req.Ids {
		in = append(in, Int64(int64(id)))
	}
	result, err := Sources.DELETE().WHERE(Sources.ID.IN(in...)).ExecContext(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to delete sources: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return &clawv1.DeleteSourceResponse{Success: rowsAffected > 0}, nil
}
