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

// UpdateImage updates an existing image
func (s *Claw) UpdateImage(ctx context.Context, req *clawv1.UpdateImageRequest) (*clawv1.UpdateImageResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	nowMillis := types.UnixMilliNow()

	columns := make([]Column, 0, 5)
	if req.PostAuthor != nil {
		columns = append(columns, Images.PostAuthor)
	}
	if req.PostAuthorUrl != nil {
		columns = append(columns, Images.PostAuthorURL)
	}
	if req.PostUrl != nil {
		columns = append(columns, Images.PostURL)
	}
	if req.IsFavorite != nil {
		columns = append(columns, Images.IsFavorite)
	}
	if len(columns) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	var imageRow model.Images
	err = Images.UPDATE(ColumnList(columns)).MODEL(model.Images{
		ID:            new(int64),
		PostAuthor:    Deref(req.PostAuthor),
		PostAuthorURL: Deref(req.PostAuthorUrl),
		PostURL:       Deref(req.PostUrl),
		IsFavorite:    types.NewBoolFromPointer(req.IsFavorite),
		UpdatedAt:     nowMillis,
	}).
		WHERE(Images.ID.EQ(Int64(req.Id))).
		RETURNING(Images.AllColumns).
		QueryContext(ctx, tx, &imageRow)
	if err != nil {
		return nil, fmt.Errorf("failed to update image: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &clawv1.UpdateImageResponse{
		Image: imageModelToProto(imageRow),
	}, nil
}
