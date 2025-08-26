package claw

import (
	"context"
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
)

// UnsubscribeDevice unsubscribes a device from sources
func (s *Claw) UnsubscribeDevice(ctx context.Context, req *clawv1.UnsubscribeDeviceRequest) (*clawv1.UnsubscribeDeviceResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if device exists
	var deviceRow model.Devices
	deviceStmt := SELECT(Devices.AllColumns).
		FROM(Devices).
		WHERE(Devices.ID.EQ(Int64(req.DeviceId)))

	err = deviceStmt.QueryContext(ctx, tx, &deviceRow)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	// Convert source IDs to expressions for the IN clause
	var sourceIDExprs []Expression
	for _, sourceID := range req.SourceIds {
		sourceIDExprs = append(sourceIDExprs, Int64(sourceID))
	}

	// Delete subscriptions
	if len(sourceIDExprs) > 0 {
		deleteStmt := DeviceSources.DELETE().
			WHERE(DeviceSources.DeviceID.EQ(Int64(req.DeviceId)).
				AND(DeviceSources.SourceID.IN(sourceIDExprs...)))

		_, err = deleteStmt.ExecContext(ctx, tx)
		if err != nil {
			return nil, fmt.Errorf("failed to delete device subscriptions: %w", err)
		}
	}

	// Get remaining subscriptions for the response
	remainingStmt := SELECT(DeviceSources.SourceID).
		FROM(DeviceSources).
		WHERE(DeviceSources.DeviceID.EQ(Int64(req.DeviceId)))

	var remainingSubscriptionRows []model.DeviceSources
	err = remainingStmt.QueryContext(ctx, tx, &remainingSubscriptionRows)
	if err != nil {
		return nil, fmt.Errorf("failed to get remaining subscriptions: %w", err)
	}

	// Convert subscriptions to slice of int64
	var subscriptions []int64
	for _, sub := range remainingSubscriptionRows {
		subscriptions = append(subscriptions, sub.SourceID)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
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

	return &clawv1.UnsubscribeDeviceResponse{
		Device: device,
	}, nil
}

