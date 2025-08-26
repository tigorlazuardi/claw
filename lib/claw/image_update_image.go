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

	// Update device assignments if provided
	if req.DeviceIds != nil {
		// Delete existing device assignments
		_, err = ImageDevices.DELETE().
			WHERE(ImageDevices.ImageID.EQ(Int64(req.Id))).
			ExecContext(ctx, tx)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing device assignments: %w", err)
		}

		// Insert new device assignments
		if len(req.DeviceIds) > 0 {
			var deviceAssignments []model.ImageDevices
			for _, deviceID := range req.DeviceIds {
				deviceAssignments = append(deviceAssignments, model.ImageDevices{
					ImageID:   req.Id,
					DeviceID:  deviceID,
					CreatedAt: nowMillis,
				})
			}

			_, err = ImageDevices.
				INSERT(ImageDevices.ImageID, ImageDevices.DeviceID, ImageDevices.CreatedAt).
				MODELS(deviceAssignments).
				ExecContext(ctx, tx)
			if err != nil {
				return nil, fmt.Errorf("failed to create device assignments: %w", err)
			}
		}
	}

	// Update tags if provided
	if req.Tags != nil {
		// Delete existing tags
		_, err = ImageTags.DELETE().
			WHERE(ImageTags.ImageID.EQ(Int64(req.Id))).
			ExecContext(ctx, tx)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing tags: %w", err)
		}

		// Insert new tags
		if len(req.Tags) > 0 {
			var imageTags []model.ImageTags
			for _, tag := range req.Tags {
				imageTags = append(imageTags, model.ImageTags{
					ImageID:   req.Id,
					Tag:       tag,
					CreatedAt: nowMillis,
				})
			}

			_, err = ImageTags.
				INSERT(ImageTags.ImageID, ImageTags.Tag, ImageTags.CreatedAt).
				MODELS(imageTags).
				ExecContext(ctx, tx)
			if err != nil {
				return nil, fmt.Errorf("failed to create image tags: %w", err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Get updated image with all related data
	imageResponse, err := s.GetImage(ctx, &clawv1.GetImageRequest{Id: req.Id})
	if err != nil {
		return nil, fmt.Errorf("failed to get updated image: %w", err)
	}

	return &clawv1.UpdateImageResponse{
		Image: imageResponse.Image,
	}, nil
}
