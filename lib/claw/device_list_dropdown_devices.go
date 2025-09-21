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
	stmt := SELECT(Devices.ID, Devices.Name).
		FROM(Devices).
		ORDER_BY(Devices.Name.ASC())

	var deviceRows []model.Devices
	err := stmt.QueryContext(ctx, s.db, &deviceRows)
	if err != nil {
		return nil, fmt.Errorf("failed to list dropdown devices: %w", err)
	}

	// Convert to protobuf dropdown options
	var items []*clawv1.ListDropDownDevicesResponse_Item
	for _, deviceRow := range deviceRows {
		items = append(items, &clawv1.ListDropDownDevicesResponse_Item{
			Id:   int32(*deviceRow.ID),
			Name: deviceRow.Name,
		})
	}

	return &clawv1.ListDropDownDevicesResponse{
		Items: items,
	}, nil
}
