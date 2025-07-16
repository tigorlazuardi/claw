package claw

import (
	"context"
	"fmt"

	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/claw/gen/table"
	"github.com/tigorlazuardi/claw/lib/claw/types"
	"github.com/go-jet/jet/v2/sqlite"
)

// UpdateSource updates an existing source
func (s *SourceService) UpdateSource(ctx context.Context, req *clawv1.UpdateSourceRequest) (*clawv1.UpdateSourceResponse, error) {
	nowMillis := types.UnixMilliNow()

	// Build dynamic update statement
	updateStmt := table.Sources.UPDATE(table.Sources.UpdatedAt).
		SET(nowMillis).
		WHERE(table.Sources.ID.EQ(sqlite.Int64(req.Id)))

	if req.Kind != nil {
		updateStmt = updateStmt.SET(table.Sources.Kind.SET(sqlite.String(*req.Kind)))
	}
	if req.Slug != nil {
		updateStmt = updateStmt.SET(table.Sources.Slug.SET(sqlite.String(*req.Slug)))
	}
	if req.DisplayName != nil {
		updateStmt = updateStmt.SET(table.Sources.DisplayName.SET(sqlite.String(*req.DisplayName)))
	}
	if req.Parameter != nil {
		updateStmt = updateStmt.SET(table.Sources.Parameter.SET(sqlite.String(*req.Parameter)))
	}
	if req.Countback != nil {
		updateStmt = updateStmt.SET(table.Sources.Countback.SET(sqlite.Int32(*req.Countback)))
	}
	if req.IsDisabled != nil {
		var disabledVal int64
		if *req.IsDisabled {
			disabledVal = 1
		} else {
			disabledVal = 0
		}
		updateStmt = updateStmt.SET(table.Sources.IsDisabled.SET(sqlite.Int64(disabledVal)))
	}

	// Execute update
	result, err := updateStmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("failed to update source: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("source not found")
	}

	// Get updated source
	getResp, err := s.GetSource(ctx, &clawv1.GetSourceRequest{Id: req.Id})
	if err != nil {
		return nil, fmt.Errorf("failed to get updated source: %w", err)
	}

	return &clawv1.UpdateSourceResponse{Source: getResp.Source}, nil
}