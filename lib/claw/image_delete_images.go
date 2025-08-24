package claw

import (
	"context"
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
)

// DeleteImages deletes images by their IDs
func (s *Claw) DeleteImages(ctx context.Context, req *clawv1.DeleteImagesRequest) (*clawv1.DeleteImagesResponse, error) {
	if len(req.Ids) == 0 {
		return &clawv1.DeleteImagesResponse{
			DeletedCount: 0,
		}, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Convert IDs to expressions
	var idExprs []Expression
	for _, id := range req.Ids {
		idExprs = append(idExprs, Int64(id))
	}

	// Delete related data first (due to foreign key constraints)
	// Delete image tags
	_, err = ImageTags.DELETE().
		WHERE(ImageTags.ImageID.IN(idExprs...)).
		ExecContext(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to delete image tags: %w", err)
	}

	// Delete image paths
	_, err = ImagePaths.DELETE().
		WHERE(ImagePaths.ImageID.IN(idExprs...)).
		ExecContext(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to delete image paths: %w", err)
	}

	// Delete image device assignments
	_, err = ImageDevices.DELETE().
		WHERE(ImageDevices.ImageID.IN(idExprs...)).
		ExecContext(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to delete image device assignments: %w", err)
	}

	// Delete images
	result, err := Images.DELETE().
		WHERE(Images.ID.IN(idExprs...)).
		ExecContext(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to delete images: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return &clawv1.DeleteImagesResponse{
		DeletedCount: int32(rowsAffected),
	}, nil
}