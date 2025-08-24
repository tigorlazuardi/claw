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

// SubscribeDevice subscribes a device to sources
func (s *Claw) SubscribeDevice(ctx context.Context, req *clawv1.SubscribeDeviceRequest) (*clawv1.SubscribeDeviceResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	nowMillis := types.UnixMilliNow()

	// Check if device exists
	var deviceRow model.Devices
	deviceStmt := SELECT(Devices.AllColumns).
		FROM(Devices).
		WHERE(Devices.ID.EQ(Int64(req.DeviceId)))

	err = deviceStmt.QueryContext(ctx, tx, &deviceRow)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	// Get existing subscriptions to avoid duplicates
	existingStmt := SELECT(DeviceSources.SourceID).
		FROM(DeviceSources).
		WHERE(DeviceSources.DeviceID.EQ(Int64(req.DeviceId)))

	var existingSubscriptions []model.DeviceSources
	err = existingStmt.QueryContext(ctx, tx, &existingSubscriptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing subscriptions: %w", err)
	}

	// Create a map of existing subscriptions for quick lookup
	existingMap := make(map[int64]bool)
	for _, sub := range existingSubscriptions {
		existingMap[sub.SourceID] = true
	}

	// Prepare new subscriptions (only those that don't already exist)
	var newSubscriptions []model.DeviceSources
	for _, sourceID := range req.SourceIds {
		if !existingMap[sourceID] {
			newSubscriptions = append(newSubscriptions, model.DeviceSources{
				DeviceID:  req.DeviceId,
				SourceID:  sourceID,
				CreatedAt: nowMillis,
			})
		}
	}

	// Insert new subscriptions if any
	if len(newSubscriptions) > 0 {
		subscriptionStmt := DeviceSources.
			INSERT(DeviceSources.DeviceID, DeviceSources.SourceID, DeviceSources.CreatedAt).
			MODELS(newSubscriptions)

		_, err = subscriptionStmt.ExecContext(ctx, tx)
		if err != nil {
			return nil, fmt.Errorf("failed to create device subscriptions: %w", err)
		}
	}

	// Get all current subscriptions for the response
	allSubscriptionsStmt := SELECT(DeviceSources.SourceID).
		FROM(DeviceSources).
		WHERE(DeviceSources.DeviceID.EQ(Int64(req.DeviceId)))

	var allSubscriptionRows []model.DeviceSources
	err = allSubscriptionsStmt.QueryContext(ctx, tx, &allSubscriptionRows)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated subscriptions: %w", err)
	}

	// Convert subscriptions to slice of int64
	var subscriptions []int64
	for _, sub := range allSubscriptionRows {
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
		ImageMinFilesize:      uint64(deviceRow.ImageMinFileSize),
		ImageMaxFilesize:      uint64(deviceRow.ImageMaxFileSize),
		Nsfw:                  clawv1.NSFWMode(deviceRow.NsfwMode),
		CreatedAt:             deviceRow.CreatedAt.ToProto(),
		UpdatedAt:             deviceRow.UpdatedAt.ToProto(),
		Subscriptions:         subscriptions,
	}

	return &clawv1.SubscribeDeviceResponse{
		Device: device,
	}, nil
}