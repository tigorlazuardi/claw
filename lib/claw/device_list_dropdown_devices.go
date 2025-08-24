package claw

import (
	"context"
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
)

// ListDropDownDevices returns a simple list of devices for dropdown selection
func (s *Claw) ListDropDownDevices(ctx context.Context, req *clawv1.ListDropDownDevicesRequest) (*clawv1.ListDropDownDevicesResponse, error) {
	// Select only slug and name for dropdown
	stmt := SELECT(Devices.Slug, Devices.Name).
		FROM(Devices).
		ORDER_BY(Devices.Name.ASC(), Devices.Slug.ASC())

	var deviceRows []model.Devices
	err := stmt.QueryContext(ctx, s.db, &deviceRows)
	if err != nil {
		return nil, fmt.Errorf("failed to list dropdown devices: %w", err)
	}

	// Convert to protobuf dropdown options
	var devices []*clawv1.DeviceDropDownOption
	for _, deviceRow := range deviceRows {
		option := &clawv1.DeviceDropDownOption{
			Slug: deviceRow.Slug,
		}

		// Set name if it's not empty
		if deviceRow.Name != "" {
			option.Name = &deviceRow.Name
		}

		devices = append(devices, option)
	}

	return &clawv1.ListDropDownDevicesResponse{
		Devices: devices,
	}, nil
}

