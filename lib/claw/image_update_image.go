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

	// Build update statement dynamically based on provided fields
	updateStmt := Images.UPDATE().SET(Images.UpdatedAt.SET(nowMillis.AsSqlLiteral()))

	// Update post author if provided
	if req.PostAuthor != nil {
		updateStmt = updateStmt.SET(Images.PostAuthor.SET(String(*req.PostAuthor)))
	}

	// Update post author URL if provided
	if req.PostAuthorUrl != nil {
		updateStmt = updateStmt.SET(Images.PostAuthorURL.SET(String(*req.PostAuthorUrl)))
	}

	// Update post URL if provided
	if req.PostUrl != nil {
		updateStmt = updateStmt.SET(Images.PostURL.SET(String(*req.PostUrl)))
	}

	// Update favorite status if provided
	if req.IsFavorite != nil {
		favoriteValue := 0
		if *req.IsFavorite {
			favoriteValue = 1
		}
		updateStmt = updateStmt.SET(Images.IsFavorite.SET(Int32(int32(favoriteValue))))
	}

	// Execute update
	finalStmt := updateStmt.WHERE(Images.ID.EQ(Int64(req.Id))).RETURNING(Images.AllColumns)

	var imageRow model.Images
	err = finalStmt.QueryContext(ctx, tx, &imageRow)
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