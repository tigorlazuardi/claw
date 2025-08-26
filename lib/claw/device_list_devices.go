package claw

import (
	"context"
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
)

// ListDevices lists devices with optional filtering and pagination
func (s *Claw) ListDevices(ctx context.Context, req *clawv1.ListDevicesRequest) (*clawv1.ListDevicesResponse, error) {
	// Build base query
	stmt := SELECT(Devices.AllColumns).FROM(Devices)

	// Add search filter if provided
	if req.Search != nil && *req.Search != "" {
		searchTerm := "%" + *req.Search + "%"
		stmt = stmt.WHERE(
			Devices.Slug.LIKE(String(searchTerm)).
				OR(Devices.Name.LIKE(String(searchTerm))),
		)
	}

	// Add pagination
	pageSize := int64(20) // default
	if req.PageSize != nil && *req.PageSize > 0 {
		pageSize = int64(*req.PageSize)
	}

	offset := int64(0)
	if req.PageToken != nil && *req.PageToken > 0 {
		offset = int64(*req.PageToken) * pageSize
	}

	stmt = stmt.ORDER_BY(Devices.CreatedAt.DESC()).LIMIT(pageSize).OFFSET(offset)

	// Execute query
	var deviceRows []model.Devices
	err := stmt.QueryContext(ctx, s.db, &deviceRows)
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}

	// Get subscriptions for all devices
	deviceIDs := make([]Expression, len(deviceRows))
	for i, deviceRow := range deviceRows {
		deviceIDs[i] = Int64(*deviceRow.ID)
	}

	var allSubscriptions []model.DeviceSources
	if len(deviceIDs) > 0 {
		subscriptionStmt := SELECT(DeviceSources.DeviceID, DeviceSources.SourceID).
			FROM(DeviceSources).
			WHERE(DeviceSources.DeviceID.IN(deviceIDs...))

		err = subscriptionStmt.QueryContext(ctx, s.db, &allSubscriptions)
		if err != nil {
			return nil, fmt.Errorf("failed to get device subscriptions: %w", err)
		}
	}

	// Group subscriptions by device ID
	subscriptionMap := make(map[int64][]int64)
	for _, sub := range allSubscriptions {
		subscriptionMap[sub.DeviceID] = append(subscriptionMap[sub.DeviceID], sub.SourceID)
	}

	// Convert to protobuf
	var devices []*clawv1.Device
	for _, deviceRow := range deviceRows {
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
			Subscriptions:         subscriptionMap[*deviceRow.ID],
		}
		devices = append(devices, device)
	}

	// Calculate next page token
	var nextPageToken *uint32
	if len(devices) == int(pageSize) {
		nextToken := uint32(offset/pageSize + 1)
		nextPageToken = &nextToken
	}

	// Get total count if search is used (optional optimization)
	var totalCount *int64
	if req.Search != nil && *req.Search != "" {
		countStmt := SELECT(COUNT(STAR)).FROM(Devices)
		if req.Search != nil && *req.Search != "" {
			searchTerm := "%" + *req.Search + "%"
			countStmt = countStmt.WHERE(
				Devices.Slug.LIKE(String(searchTerm)).
					OR(Devices.Name.LIKE(String(searchTerm))),
			)
		}

		var count int64
		err = countStmt.QueryContext(ctx, s.db, &count)
		if err == nil {
			totalCount = &count
		}
	}

	return &clawv1.ListDevicesResponse{
		Devices:       devices,
		NextPageToken: nextPageToken,
		TotalCount:    totalCount,
	}, nil
}
