package claw

import (
	"context"
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/otel"
)

// ListSources lists sources with optional filtering and cursor-based pagination
func (s *Claw) ListSources(ctx context.Context, req *clawv1.ListSourcesRequest) (*clawv1.ListSourcesResponse, error) {
	cond := Bool(true)

	if req.Name != nil {
		cond.AND(Sources.Name.EQ(String(*req.Name)))
	}
	if search := req.GetSearch(); search != "" {
		searchTerm := String("%" + search + "%")
		cond.AND(
			Sources.Name.LIKE(searchTerm).
				OR(Sources.DisplayName.LIKE(searchTerm)),
		)
	}

	// Handle cursor pagination
	if req.Pagination != nil {
		if token := req.Pagination.GetNextToken(); token != 0 {
			cond.AND(Sources.ID.GT(Int64(int64(token))))
		}
		if token := req.Pagination.GetPrevToken(); token != 0 {
			cond.AND(Sources.ID.LT(Int64(int64(token))))
		}
	}

	sorts := make([]OrderByClause, 0, len(req.Sorts)+1)
	for _, sort := range req.Sorts {
		var col OrderByClause
		switch sort.Field {
		case clawv1.SourceSortField_SOURCE_SORT_FIELD_NAME:
			col = toOrderByClause(Sources.Name, sort.Desc)
		case clawv1.SourceSortField_SOURCE_SORT_FIELD_DISPLAY_NAME:
			col = toOrderByClause(Sources.DisplayName, sort.Desc)
		case clawv1.SourceSortField_SOURCE_SORT_FIELD_COUNTBACK:
			col = toOrderByClause(Sources.Countback, sort.Desc)
		case clawv1.SourceSortField_SOURCE_SORT_FIELD_IS_DISABLED:
			col = toOrderByClause(Sources.IsDisabled, sort.Desc)
		case clawv1.SourceSortField_SOURCE_SORT_FIELD_LAST_RUN_AT:
			col = toOrderByClause(Sources.LastRunAt, sort.Desc)
		case clawv1.SourceSortField_SOURCE_SORT_FIELD_CREATED_AT:
			col = toOrderByClause(Sources.CreatedAt, sort.Desc)
		case clawv1.SourceSortField_SOURCE_SORT_FIELD_UPDATED_AT:
			col = toOrderByClause(Sources.UpdatedAt, sort.Desc)
		default:
			continue
		}
		sorts = append(sorts, col)
	}
	sorts = append(sorts, Sources.ID.ASC()) // Always add ID as the last sort for consistent ordering

	limit := uint32(25)
	if req.Pagination != nil && req.Pagination.GetSize() != 0 {
		limit = Clamp(req.Pagination.GetSize(), 1, 100)
	}

	var (
		from         ReadableTable = Sources
		extraColumns []Projection
	)
	if req.GetIncludeSchedules() {
		extraColumns = append(extraColumns, Schedules.AllColumns)
		from = Sources.LEFT_JOIN(Schedules, Sources.ID.EQ(Schedules.SourceID))
	}

	query := SELECT(Sources.AllColumns, extraColumns...).
		FROM(from).
		WHERE(cond).
		ORDER_BY(sorts...).
		LIMIT(int64(limit))

	var rows []struct {
		model.Sources
		Schedules []model.Schedules
	}
	ctx = otel.ContextWithDatabaseCaller(ctx, otel.CurrentCaller())
	err := query.QueryContext(ctx, s.db, &rows)
	if err != nil {
		return nil, fmt.Errorf("failed to list sources: %w", err)
	}
	if len(rows) == 0 {
		return &clawv1.ListSourcesResponse{
			Items: []*clawv1.ListSourcesResponse_Item{},
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

	items := make([]*clawv1.ListSourcesResponse_Item, 0, len(rows))
	for _, row := range rows {
		item := &clawv1.ListSourcesResponse_Item{
			Schedules: make([]*clawv1.SourceSchedule, 0, len(row.Schedules)),
		}
		item.Source = &clawv1.Source{
			Id:          int64(*row.ID),
			Name:        row.Name,
			DisplayName: row.DisplayName,
			Parameter:   row.Parameter,
			Countback:   int32(row.Countback),
			IsDisabled:  row.IsDisabled.Bool(),
			CreatedAt:   row.CreatedAt.ToProto(),
			UpdatedAt:   row.UpdatedAt.ToProto(),
		}
		for _, sched := range row.Schedules {
			item.Schedules = append(item.Schedules, &clawv1.SourceSchedule{
				Id:        int64(*sched.ID),
				Schedule:  sched.Schedule,
				CreatedAt: sched.CreatedAt.ToProto(),
			})
		}
		items = append(items, item)
	}

	response := &clawv1.ListSourcesResponse{
		Items:      items,
		Pagination: pagination,
	}

	return response, nil
}
