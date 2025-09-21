package claw

import (
	"context"
	"fmt"

	"github.com/go-jet/jet/v2/sqlite"
	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/claw/types"
)

// UpdateDevice updates an existing device
func (s *Claw) UpdateDevice(ctx context.Context, req *clawv1.UpdateDeviceRequest) (*clawv1.UpdateDeviceResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	columns := make([]Column, 0, 10)
	if req.Name != nil {
		columns = append(columns, Devices.Name)
	}
	if req.Height != nil {
		columns = append(columns, Devices.Height)
	}
	if req.Width != nil {
		columns = append(columns, Devices.Width)
	}
	if req.AspectRatioDifference != nil {
		columns = append(columns, Devices.AspectRatioDifference)
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
	if req.Nsfw != nil {
		columns = append(columns, Devices.NsfwMode)
	}
	if len(columns) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	var out model.Devices
	err = Devices.UPDATE(ColumnList(columns)).MODEL(model.Devices{
		Name:                  Deref(req.Name),
		Width:                 int64(Deref(req.Width)),
		Height:                int64(Deref(req.Height)),
		FilenameTemplate:      Deref(req.FilenameTemplate),
		ImageMaxWidth:         int64(Deref(req.ImageMaxWidth)),
		ImageMinHeight:        int64(Deref(req.ImageMinHeight)),
		ImageMaxHeight:        int64(Deref(req.ImageMaxHeight)),
		AspectRatioDifference: Deref(req.AspectRatioDifference),
		NsfwMode:              int64(*req.Nsfw),
		ImageMinFileSize:      int64(Deref(req.ImageMinFilesize)),
		ImageMaxFileSize:      int64(Deref(req.ImageMaxFilesize)),
		ImageMinWidth:         int64(Deref(req.ImageMinWidth)),
		UpdatedAt:             types.UnixMilliNow(),
	}).
		WHERE(Devices.ID.EQ(sqlite.Int(int64(req.Id)))).
		RETURNING(Devices.AllColumns).
		QueryContext(ctx, s.db, &out)
	if err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}

	// Convert to protobuf
	device := &clawv1.Device{
		Id:                    *out.ID,
		Name:                  out.Name,
		Height:                int32(out.Height),
		Width:                 int32(out.Width),
		AspectRatioDifference: out.AspectRatioDifference,
		FilenameTemplate:      &out.FilenameTemplate,
		ImageMinHeight:        uint32(out.ImageMinHeight),
		ImageMinWidth:         uint32(out.ImageMinWidth),
		ImageMaxHeight:        uint32(out.ImageMaxHeight),
		ImageMaxWidth:         uint32(out.ImageMaxWidth),
		ImageMinFilesize:      uint32(out.ImageMinFileSize),
		ImageMaxFilesize:      uint32(out.ImageMaxFileSize),
		Nsfw:                  clawv1.NSFWMode(out.NsfwMode),
		CreatedAt:             out.CreatedAt.ToProto(),
		UpdatedAt:             out.UpdatedAt.ToProto(),
	}

	return &clawv1.UpdateDeviceResponse{
		Device: device,
	}, nil
}
