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

// MarkFavorite marks or unmarks images as favorite
func (s *Claw) MarkFavorite(ctx context.Context, req *clawv1.MarkFavoriteRequest) (*clawv1.MarkFavoriteResponse, error) {
	if len(req.ImageIds) == 0 {
		return &clawv1.MarkFavoriteResponse{
			UpdatedCount: 0,
		}, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	nowMillis := types.UnixMilliNow()

	// Convert IDs to expressions
	var idExprs []Expression
	for _, id := range req.ImageIds {
		idExprs = append(idExprs, Int64(id))
	}
	result, err := Images.UPDATE(Images.IsFavorite, Images.UpdatedAt).
		MODEL(model.Images{
			IsFavorite: types.NewBool(req.IsFavorite),
			UpdatedAt:  nowMillis,
		}).
		WHERE(Images.ID.IN(idExprs...)).
		ExecContext(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to update favorite status: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return &clawv1.MarkFavoriteResponse{
		UpdatedCount: int32(rowsAffected),
	}, nil
}

