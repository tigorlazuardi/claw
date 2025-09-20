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
	cond := Bool(true)

	if search := req.GetSearch(); search != "" {
		searchTerm := String("%" + search + "%")
		cond.AND(
			Sources.Name.LIKE(searchTerm).OR(Sources.DisplayName.LIKE(searchTerm)),
		)
	}

	// Add pagination
	limit := uint32(20) // default
	if req.Pagination != nil {
		if size := req.Pagination.GetSize(); size > 0 {
			limit = Clamp(size, 1, 100)
		}
		if token := req.Pagination.GetNextToken(); token > 0 {
			cond = cond.AND(Devices.ID.GT(Int64(int64(token))))
		}
		if token := req.Pagination.GetPrevToken(); token > 0 {
			cond = cond.AND(Devices.ID.LT(Int64(int64(token))))
		}
	}

	sorts := make([]OrderByClause, 0, len(req.Sorts)+1)
	for _, sort := range req.Sorts {
		var col OrderByClause
		switch sort.Field {
		case clawv1.DeviceSortField_DEVICE_SORT_FIELD_NAME:
			col = toOrderByClause(Devices.Name, sort.Desc)
		case clawv1.DeviceSortField_DEVICE_SORT_FIELD_HEIGHT:
			col = toOrderByClause(Devices.Height, sort.Desc)
		case clawv1.DeviceSortField_DEVICE_SORT_FIELD_WIDTH:
			col = toOrderByClause(Devices.Width, sort.Desc)
		case clawv1.DeviceSortField_DEVICE_SORT_FIELD_ASPECT_RATIO_DIFFERENCE:
			col = toOrderByClause(Devices.AspectRatioDifference, sort.Desc)
		case clawv1.DeviceSortField_DEVICE_SORT_FIELD_NSFW:
			col = toOrderByClause(Devices.NsfwMode, sort.Desc)
		case clawv1.DeviceSortField_DEVICE_SORT_FIELD_CREATED_AT:
			col = toOrderByClause(Devices.CreatedAt, sort.Desc)
		case clawv1.DeviceSortField_DEVICE_SORT_FIELD_UPDATED_AT:
			col = toOrderByClause(Devices.UpdatedAt, sort.Desc)
		default:
			continue
		}
		sorts = append(sorts, col)
	}
	if len(sorts) == 0 {
		sorts = append(sorts,
			Devices.LastActiveAt.DESC(),
			Devices.IsDisabled.DESC(),
			Devices.Name.ASC(),
		) // default sort
	}
	sorts = append(sorts, Devices.ID.ASC()) // tie-breaker

	var (
		from         ReadableTable = Devices
		extraColumns []Projection
		groupBy      GroupByClause
	)
	if req.GetCountImages() {
		from = from.INNER_JOIN(ImageDevices, ImageDevices.DeviceID.EQ(Devices.ID))
		extraColumns = append(extraColumns, COUNT(ImageDevices.ID).AS("image_count"))
		groupBy = Devices.ID
	}
	if sourceID := req.GetSourceId(); sourceID != 0 {
		from = from.INNER_JOIN(DeviceSources, DeviceSources.DeviceID.EQ(Devices.ID))
		cond = cond.AND(DeviceSources.SourceID.EQ(Int64(int64(sourceID))))
	}
	stmt := SELECT(Devices.AllColumns, extraColumns...).
		FROM(from).
		WHERE(cond).
		ORDER_BY(sorts...).
		LIMIT(int64(limit))
	if groupBy != nil {
		stmt = stmt.GROUP_BY(groupBy)
	}

	var rows []struct {
		model.Devices
		ImageCount *int64
	}
	err := stmt.QueryContext(ctx, s.db, &rows)
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}
	if len(rows) == 0 {
		return &clawv1.ListDevicesResponse{
			Items: []*clawv1.ListDevicesResponse_Item{},
			Pagination: &clawv1.Pagination{
				Size: &limit,
			},
		}, nil
	}

	pagination := &clawv1.Pagination{
		Size: &limit,
	}

	if len(rows) >= int(limit) {
		lastRow := rows[len(rows)-1]
		pagination.NextToken = Ptr(uint32(*lastRow.ID))
	}
	if req.Pagination != nil && req.Pagination.GetPrevToken() > 0 { // Not first page
		firstRow := rows[0]
		pagination.PrevToken = Ptr(uint32(*firstRow.ID))
	}

	items := make([]*clawv1.ListDevicesResponse_Item, 0, len(rows))
	for _, row := range rows {
		item := &clawv1.ListDevicesResponse_Item{
			Device: &clawv1.Device{
				Id:                    int64(*row.ID),
				Name:                  Ptr(row.Name),
				Slug:                  row.Slug,
				Width:                 int32(row.Width),
				Height:                int32(row.Height),
				AspectRatioDifference: row.AspectRatioDifference,
				Nsfw:                  clawv1.NSFWMode(int32(row.NsfwMode)),
				CreatedAt:             row.CreatedAt.ToProto(),
				UpdatedAt:             row.UpdatedAt.ToProto(),
			},
			ImageCount: row.ImageCount,
		}
		items = append(items, item)
	}

	return &clawv1.ListDevicesResponse{
		Items:      items,
		Pagination: pagination,
	}, nil
}
