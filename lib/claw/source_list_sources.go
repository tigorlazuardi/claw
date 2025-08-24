package claw

import (
	"context"
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
)

// ListSources lists sources with optional filtering and cursor-based pagination
func (s *Claw) ListSources(ctx context.Context, req *clawv1.ListSourcesRequest) (*clawv1.ListSourcesResponse, error) {
	cond := Bool(true)

	if req.Kind != nil {
		cond.AND(Sources.Kind.EQ(String(*req.Kind)))
	}
	if search := req.GetSearch(); search != "" {
		searchTerm := String("%" + search + "%")
		cond.AND(
			Sources.Kind.LIKE(searchTerm).
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
		case clawv1.SourceSortField_SOURCE_SORT_FIELD_KIND:
			col = toOrderByClause(Sources.Kind, sort.Desc)
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

	var from ReadableTable = Sources
	if req.GetIncludeSchedules() {
		from = Sources.LEFT_JOIN(Schedules, Sources.ID.EQ(Schedules.SourceID))
	}

	query := SELECT(Sources.AllColumns).
		FROM(from).
		WHERE(cond).
		ORDER_BY(sorts...).
		LIMIT(int64(limit))

	var rows []struct {
		model.Sources
		Schedules []model.Schedules
	}
	err := query.QueryContext(ctx, s.db, &rows)
	if err != nil {
		return nil, fmt.Errorf("failed to list sources: %w", err)
	}
	if len(rows) == 0 {
		return &clawv1.ListSourcesResponse{
			Sources: []*clawv1.Source{},
			Pagination: &clawv1.Pagination{
				Size: &limit,
			},
		}, nil
	}

	// Check if there are more results and set next page token
	var nextPageToken, prevPageToken uint32

	if len(rows) >= int(limit) {
		firstRow, lastRow := rows[0], rows[len(rows)-1]
		nextPageToken = uint32(*lastRow.ID)
		prevPageToken = uint32(*firstRow.ID)
	}
	// Convert to protobuf
	sources := make([]*clawv1.Source, 0, len(rows))
	for _, row := range rows {
		source := &clawv1.Source{
			Id:          int64(*row.ID),
			Kind:        row.Kind,
			DisplayName: row.DisplayName,
			Parameter:   row.Parameter,
			Countback:   int32(row.Countback),
			IsDisabled:  row.IsDisabled.Bool(),
			CreatedAt:   row.CreatedAt.ToProto(),
			UpdatedAt:   row.UpdatedAt.ToProto(),
			Schedules:   make([]*clawv1.SourceSchedule, 0, len(row.Schedules)),
		}
		for _, sched := range row.Schedules {
			source.Schedules = append(source.Schedules, &clawv1.SourceSchedule{
				Id:        int64(*sched.ID),
				Schedule:  sched.Schedule,
				CreatedAt: sched.CreatedAt.ToProto(),
			})
		}
		sources = append(sources, source)
	}

	response := &clawv1.ListSourcesResponse{
		Sources: sources,
		Pagination: &clawv1.Pagination{
			Size:      &limit,
			NextToken: &nextPageToken,
			PrevToken: &prevPageToken,
		},
	}

	return response, nil
}
