package claw

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/sqlite"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/claw/gen/table"
)

// DeleteSource deletes a source by ID
func (s *Claw) DeleteSource(ctx context.Context, req *clawv1.DeleteSourceRequest) (*clawv1.DeleteSourceResponse, error) {
	deleteStmt := table.Sources.DELETE().WHERE(table.Sources.ID.EQ(sqlite.Int64(req.Id)))

	result, err := deleteStmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to delete source: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return &clawv1.DeleteSourceResponse{Success: rowsAffected > 0}, nil
}

