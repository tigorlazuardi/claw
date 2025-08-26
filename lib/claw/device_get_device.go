package claw

import (
	"context"
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
)

// GetDevice retrieves a device by ID
func (s *Claw) GetDevice(ctx context.Context, req *clawv1.GetDeviceRequest) (*clawv1.GetDeviceResponse, error) {
	stmt := SELECT(Devices.AllColumns).
		FROM(Devices).
		WHERE(Devices.ID.EQ(Int64(req.Id)))

	var deviceRow model.Devices

	err := stmt.QueryContext(ctx, s.db, &deviceRow)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	// Get device subscriptions
	subscriptionStmt := SELECT(DeviceSources.SourceID).
		FROM(DeviceSources).
		WHERE(DeviceSources.DeviceID.EQ(Int64(req.Id)))

	var subscriptionRows []model.DeviceSources
	err = subscriptionStmt.QueryContext(ctx, s.db, &subscriptionRows)
	if err != nil {
		return nil, fmt.Errorf("failed to get device subscriptions: %w", err)
	}

	// Convert subscriptions to slice of int64
	var subscriptions []int64
	for _, sub := range subscriptionRows {
		subscriptions = append(subscriptions, sub.SourceID)
	}

	// Convert to protobuf
	device := &clawv1.Device{
		Id:                    *deviceRow.ID,
		Slug:                  deviceRow.Slug,
		Name:                  &deviceRow.Name,
		Height:                int32(deviceRow.Height),
		Width:                 int32(deviceRow.Width),
		AspectRatioDifference: deviceRow.AspectRatioDifference,
		SaveDir:               deviceRow.SaveDir,
		FilenameTemplate:      &deviceRow.FilenameTemplate,
		ImageMinHeight:        uint32(deviceRow.ImageMinHeight),
		ImageMinWidth:         uint32(deviceRow.ImageMinWidth),
		ImageMaxHeight:        uint32(deviceRow.ImageMaxHeight),
		ImageMaxWidth:         uint32(deviceRow.ImageMaxWidth),
		ImageMinFilesize:      uint32(deviceRow.ImageMinFileSize),
		ImageMaxFilesize:      uint32(deviceRow.ImageMaxFileSize),
		Nsfw:                  clawv1.NSFWMode(deviceRow.NsfwMode),
		CreatedAt:             deviceRow.CreatedAt.ToProto(),
		UpdatedAt:             deviceRow.UpdatedAt.ToProto(),
		Subscriptions:         subscriptions,
	}

	return &clawv1.GetDeviceResponse{
		Device: device,
	}, nil
}
