package claw

import (
	"context"
	"fmt"

	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/claw/types"
)

// CreateDevice creates a new device
func (s *Claw) CreateDevice(ctx context.Context, req *clawv1.CreateDeviceRequest) (*clawv1.CreateDeviceResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	nowMillis := types.UnixMilliNow()

	// Convert proto values to model types
	name := ""
	if req.Name != nil {
		name = *req.Name
	}

	saveDir := ""
	if req.SaveDir != nil {
		saveDir = *req.SaveDir
	}

	filenameTemplate := ""
	if req.FilenameTemplate != nil {
		filenameTemplate = *req.FilenameTemplate
	}

	minHeight := int64(0)
	if req.ImageMinHeight != nil {
		minHeight = int64(*req.ImageMinHeight)
	}

	maxHeight := int64(0)
	if req.ImageMaxHeight != nil {
		maxHeight = int64(*req.ImageMaxHeight)
	}

	minWidth := int64(0)
	if req.ImageMinWidth != nil {
		minWidth = int64(*req.ImageMinWidth)
	}

	maxWidth := int64(0)
	if req.ImageMaxWidth != nil {
		maxWidth = int64(*req.ImageMaxWidth)
	}

	minFileSize := int64(0)
	if req.ImageMinFilesize != nil {
		minFileSize = int64(*req.ImageMinFilesize)
	}

	maxFileSize := int64(0)
	if req.ImageMaxFilesize != nil {
		maxFileSize = int64(*req.ImageMaxFilesize)
	}

	// Insert device
	deviceStmt := Devices.INSERT(
		Devices.Slug,
		Devices.Name,
		Devices.SaveDir,
		Devices.FilenameTemplate,
		Devices.Width,
		Devices.Height,
		Devices.AspectRatioDifference,
		Devices.ImageMinWidth,
		Devices.ImageMaxWidth,
		Devices.ImageMinHeight,
		Devices.ImageMaxHeight,
		Devices.ImageMinFileSize,
		Devices.ImageMaxFileSize,
		Devices.NsfwMode,
		Devices.CreatedAt,
		Devices.UpdatedAt,
	).MODEL(model.Devices{
		Slug:                  req.Slug,
		Name:                  name,
		SaveDir:               saveDir,
		FilenameTemplate:      filenameTemplate,
		Width:                 int64(req.Width),
		Height:                int64(req.Height),
		AspectRatioDifference: req.AspectRatioDifference,
		ImageMinWidth:         minWidth,
		ImageMaxWidth:         maxWidth,
		ImageMinHeight:        minHeight,
		ImageMaxHeight:        maxHeight,
		ImageMinFileSize:      minFileSize,
		ImageMaxFileSize:      maxFileSize,
		NsfwMode:              int64(req.Nsfw),
		CreatedAt:             nowMillis,
		UpdatedAt:             nowMillis,
	}).RETURNING(Devices.AllColumns)

	var deviceRow model.Devices

	err = deviceStmt.QueryContext(ctx, tx, &deviceRow)
	if err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	// Create device-source subscriptions if provided
	if len(req.Subscriptions) > 0 {
		var subscriptions []model.DeviceSources
		for _, sourceID := range req.Subscriptions {
			subscriptions = append(subscriptions, model.DeviceSources{
				DeviceID:  *deviceRow.ID,
				SourceID:  sourceID,
				CreatedAt: nowMillis,
			})
		}
		
		subscriptionStmt := DeviceSources.
			INSERT(DeviceSources.DeviceID, DeviceSources.SourceID, DeviceSources.CreatedAt).
			MODELS(subscriptions)

		_, err = subscriptionStmt.ExecContext(ctx, tx)
		if err != nil {
			return nil, fmt.Errorf("failed to create device subscriptions: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Convert to protobuf
	device := &clawv1.Device{
		Id:                      *deviceRow.ID,
		Slug:                    deviceRow.Slug,
		Name:                    &deviceRow.Name,
		Height:                  int32(deviceRow.Height),
		Width:                   int32(deviceRow.Width),
		AspectRatioDifference:   deviceRow.AspectRatioDifference,
		SaveDir:                 deviceRow.SaveDir,
		FilenameTemplate:        &deviceRow.FilenameTemplate,
		ImageMinHeight:          uint32(deviceRow.ImageMinHeight),
		ImageMinWidth:           uint32(deviceRow.ImageMinWidth),
		ImageMaxHeight:          uint32(deviceRow.ImageMaxHeight),
		ImageMaxWidth:           uint32(deviceRow.ImageMaxWidth),
		ImageMinFilesize:        uint64(deviceRow.ImageMinFileSize),
		ImageMaxFilesize:        uint64(deviceRow.ImageMaxFileSize),
		Nsfw:                    clawv1.NSFWMode(deviceRow.NsfwMode),
		CreatedAt:               deviceRow.CreatedAt.ToProto(),
		UpdatedAt:               deviceRow.UpdatedAt.ToProto(),
		Subscriptions:           req.Subscriptions, // Return the subscriptions that were created
	}

	return &clawv1.CreateDeviceResponse{
		Device: device,
	}, nil
}