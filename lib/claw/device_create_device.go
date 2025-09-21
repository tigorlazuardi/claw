package claw

import (
	"context"
	"database/sql"
	"fmt"
	"slices"

	. "github.com/go-jet/jet/v2/sqlite"
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

	columns := ColumnList{
		Devices.Name,
		Devices.Slug,
		Devices.Width,
		Devices.Height,
		Devices.AspectRatioDifference,
		Devices.NsfwMode,
		Devices.CreatedAt,
		Devices.UpdatedAt,
	}
	if req.FilenameTemplate != nil {
		columns = append(columns, Devices.FilenameTemplate)
	}
	if req.ImageMinHeight != nil {
		columns = append(columns, Devices.ImageMinHeight)
	}
	if req.ImageMaxHeight != nil {
		columns = append(columns, Devices.ImageMaxHeight)
	}
	if req.ImageMinWidth != nil {
		columns = append(columns, Devices.ImageMinWidth)
	}
	if req.ImageMaxWidth != nil {
		columns = append(columns, Devices.ImageMaxWidth)
	}
	if req.ImageMinFilesize != nil {
		columns = append(columns, Devices.ImageMinFileSize)
	}
	if req.ImageMaxFilesize != nil {
		columns = append(columns, Devices.ImageMaxFileSize)
	}
	if req.IsDisabled != nil {
		columns = append(columns, Devices.IsDisabled)
	}

	// Insert device
	deviceStmt := Devices.INSERT(columns).MODEL(model.Devices{
		Slug:                  req.Slug,
		Name:                  req.Name,
		FilenameTemplate:      Deref(req.FilenameTemplate),
		Width:                 int64(req.Width),
		Height:                int64(req.Height),
		AspectRatioDifference: req.AspectRatioDifference,
		ImageMinWidth:         int64(Deref(req.ImageMinWidth)),
		ImageMaxWidth:         int64(Deref(req.ImageMaxWidth)),
		ImageMinHeight:        int64(Deref(req.ImageMinHeight)),
		ImageMaxHeight:        int64(Deref(req.ImageMaxHeight)),
		ImageMinFileSize:      int64(Deref(req.ImageMinFilesize)),
		ImageMaxFileSize:      int64(Deref(req.ImageMaxFilesize)),
		NsfwMode:              int64(req.Nsfw),
		CreatedAt:             nowMillis,
		UpdatedAt:             nowMillis,
		IsDisabled:            types.Bool(Deref(req.IsDisabled)),
	}).RETURNING(Devices.AllColumns)

	var deviceRow model.Devices

	err = deviceStmt.QueryContext(ctx, tx, &deviceRow)
	if err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	// Create device-source subscriptions if provided
	if len(req.Sources) > 0 {
		if err := s.validateSubscriptionExists(ctx, tx, req.Sources); err != nil {
			return nil, err
		}
		var subscriptions []model.DeviceSources
		for _, sourceID := range req.Sources {
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
		Id:                    *deviceRow.ID,
		Name:                  deviceRow.Name,
		Height:                int32(deviceRow.Height),
		Width:                 int32(deviceRow.Width),
		AspectRatioDifference: deviceRow.AspectRatioDifference,
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
	}

	return &clawv1.CreateDeviceResponse{
		Device:  device,
		Sources: req.Sources,
	}, nil
}

func (claw *Claw) validateSubscriptionExists(ctx context.Context, tx *sql.Tx, sourceIDs []int64) error {
	if len(sourceIDs) == 0 {
		return nil
	}
	// validate source IDs
	ids := make([]Expression, len(sourceIDs))
	for _, sourceID := range sourceIDs {
		ids = append(ids, Int64(sourceID))
	}
	type SourceID struct {
		ID int64
	}
	var out []SourceID
	err := SELECT(Sources.ID).WHERE(Sources.ID.IN(ids...)).QueryContext(ctx, tx, &out)
	if err != nil {
		return fmt.Errorf("failed to validate source IDs: %w", err)
	}
	if len(out) != len(sourceIDs) {
		missingIds := make([]int64, 0, len(out))
		for i := range out {
			if !slices.ContainsFunc(out, func(db SourceID) bool {
				return db.ID == sourceIDs[i]
			}) {
				missingIds = append(missingIds, sourceIDs[i])
			}
		}
		return fmt.Errorf("one or more source IDs do not exist in database: %v", missingIds)
	}
	return nil
}
