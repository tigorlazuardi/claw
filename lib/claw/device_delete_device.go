package claw

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/sqlite"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
)

// DeleteDevice deletes devices by their slugs
func (s *Claw) DeleteDevice(ctx context.Context, req *clawv1.DeleteDeviceRequest) (*clawv1.DeleteDeviceResponse, error) {
	if len(req.Slugs) == 0 {
		return &clawv1.DeleteDeviceResponse{
			Success: true,
		}, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Convert slugs to sqlite expressions
	var slugExprs []sqlite.Expression
	for _, slug := range req.Slugs {
		slugExprs = append(slugExprs, sqlite.String(slug))
	}

	// Delete devices by slugs
	deleteStmt := Devices.DELETE().WHERE(Devices.Slug.IN(slugExprs...))

	_, err = deleteStmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to delete devices: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &clawv1.DeleteDeviceResponse{
		Success: true,
	}, nil
}