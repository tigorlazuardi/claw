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
	stmt := SELECT(
		Devices.AllColumns,
		Sources.AllColumns,
		COUNT(ImageDevices.DeviceID).AS("image_count"),
	).
		FROM(
			Devices.INNER_JOIN(DeviceSources, DeviceSources.DeviceID.EQ(Devices.ID)).
				INNER_JOIN(Sources, Sources.ID.EQ(DeviceSources.SourceID)).
				INNER_JOIN(ImageDevices, ImageDevices.DeviceID.EQ(Devices.ID)),
		).
		WHERE(Devices.ID.EQ(Int64(req.Id)))

	var row struct {
		model.Devices
		ImageCount int32
		Sources    []model.Sources
	}

	err := stmt.QueryContext(ctx, s.db, &row)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	resp := &clawv1.GetDeviceResponse{
		Device:     deviceModelToProto(row.Devices),
		ImageCount: row.ImageCount,
		Sources:    make([]*clawv1.Source, 0, len(row.Sources)),
	}

	for _, source := range row.Sources {
		resp.Sources = append(resp.Sources, sourceModelToProto(source))
	}

	return resp, nil
}
