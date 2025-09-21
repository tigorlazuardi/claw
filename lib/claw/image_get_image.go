package claw

import (
	"context"
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
)

// GetImage retrieves an image by ID with all related data
func (s *Claw) GetImage(ctx context.Context, req *clawv1.GetImageRequest) (*clawv1.GetImageResponse, error) {
	stmt := SELECT(
		Images.AllColumns,
		Tags.AllColumns,
		ImageDevices.AllColumns.As("assignments.image_devices"),
		Devices.AllColumns.As("assignments.devices"),
		Sources.AllColumns,
	).
		FROM(
			Images.INNER_JOIN(ImageTags, ImageTags.ImageID.EQ(Int64(req.Id))).
				INNER_JOIN(Tags, Tags.ID.EQ(ImageTags.TagID)).
				INNER_JOIN(ImageDevices, ImageDevices.ImageID.EQ(Int64(req.Id))).
				INNER_JOIN(Devices, Devices.ID.EQ(ImageDevices.DeviceID)).
				INNER_JOIN(Sources, Sources.ID.EQ(Images.SourceID)),
		).
		WHERE(Images.ID.EQ(Int64(req.Id)))

	var row struct {
		model.Images
		Tags        []model.Tags
		Assignments []struct {
			ImageDevices model.ImageDevices
			Devices      model.Devices
		}
		Sources model.Sources
	}

	err := stmt.QueryContext(ctx, s.db, &row)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	assignments := make([]*clawv1.GetImageResponse_Assignment, 0, len(row.Assignments))
	for _, assign := range row.Assignments {
		assignments = append(assignments, &clawv1.GetImageResponse_Assignment{
			Path:   assign.ImageDevices.Path,
			Device: deviceModelToProto(assign.Devices),
		})
	}

	return &clawv1.GetImageResponse{
		Image:       imageModelToProto(row.Images),
		Tags:        tagModelsToProto(row.Tags),
		Assignments: assignments,
		Source:      sourceModelToProto(row.Sources),
	}, nil
}
