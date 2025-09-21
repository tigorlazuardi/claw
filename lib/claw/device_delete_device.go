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
	if len(req.Ids) == 0 {
		return &clawv1.DeleteDeviceResponse{}, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Convert slugs to sqlite expressions
	var idExpr []sqlite.Expression
	for _, id := range req.Ids {
		idExpr = append(idExpr, sqlite.Int(int64(id)))
	}

	// Delete devices by slugs
	deleteStmt := Devices.DELETE().WHERE(Devices.ID.IN(idExpr...))

	_, err = deleteStmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to delete devices: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &clawv1.DeleteDeviceResponse{}, nil
}
