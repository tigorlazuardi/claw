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

	var exist struct {
		ID int
	}

	err = SELECT(Devices.ID).
		FROM(Devices).
		WHERE(Devices.ID.EQ(Int64(req.DeviceId))).
		QueryContext(ctx, tx, &exist)
	if err != nil {
		return nil, fmt.Errorf("failed to check if device exist on database: %w", err)
	}
	if exist.ID == 0 {
		return nil, fmt.Errorf("device with id %d does not exist", req.DeviceId)
	}

	insertModels := make([]model.DeviceSources, 0, len(req.SourceIds))
	for _, id := range req.SourceIds {
		insertModels = append(insertModels, model.DeviceSources{
			DeviceID:  req.DeviceId,
			SourceID:  id,
			CreatedAt: nowMillis,
		})
	}

	_, err = DeviceSources.
		INSERT(DeviceSources.DeviceID, DeviceSources.SourceID).
		MODELS(insertModels).
		ON_CONFLICT(DeviceSources.DeviceID, DeviceSources.SourceID).
		DO_NOTHING().
		ExecContext(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe device to sources: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &clawv1.SubscribeDeviceResponse{}, nil
}
