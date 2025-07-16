package claw

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-jet/jet/v2/sqlite"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/claw/gen/table"
	"github.com/tigorlazuardi/claw/lib/claw/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ListSources lists sources with optional filtering and cursor-based pagination
func (s *Claw) ListSources(ctx context.Context, req *clawv1.ListSourcesRequest) (*clawv1.ListSourcesResponse, error) {
	// Build query with optional filters
	query := sqlite.SELECT(table.Sources.AllColumns).FROM(table.Sources)

	if req.Kind != nil {
		query = query.WHERE(table.Sources.Kind.EQ(sqlite.String(*req.Kind)))
	}
	if req.Slug != nil {
		query = query.WHERE(table.Sources.Slug.EQ(sqlite.String(*req.Slug)))
	}

	// Handle cursor pagination
	if req.PageToken != "" {
		cursorID, err := strconv.ParseInt(req.PageToken, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid page token: %w", err)
		}
		query = query.WHERE(table.Sources.ID.GT(sqlite.Int64(cursorID)))
	}

	// Always sort by ID for consistent pagination and add limit + 1 to check if there's a next page
	query = query.ORDER_BY(table.Sources.ID.ASC()).LIMIT(int64(req.PageSize + 1))

	var sourceRows []struct {
		ID          int64 `sql:"primary_key"`
		Kind        string
		Slug        string
		DisplayName string
		Parameter   string
		Countback   int32
		IsDisabled  types.Bool
		LastRunAt   *types.UnixMilli
		CreatedAt   types.UnixMilli
		UpdatedAt   types.UnixMilli
	}

	err := query.QueryContext(ctx, s.db, &sourceRows)
	if err != nil {
		return nil, fmt.Errorf("failed to list sources: %w", err)
	}

	// Check if there are more results and set next page token
	var nextPageToken string
	hasMore := len(sourceRows) > int(req.PageSize)
	if hasMore {
		// Remove the extra row we fetched for pagination check
		sourceRows = sourceRows[:req.PageSize]
		// Next page token is the ID of the last item in current page
		nextPageToken = strconv.FormatInt(sourceRows[len(sourceRows)-1].ID, 10)
	}

	// Convert to protobuf
	var sources []*clawv1.SourceData
	schedulesMap := make(map[int64]*clawv1.SourceScheduleList)

	for _, row := range sourceRows {
		var lastRunAt *timestamppb.Timestamp
		if row.LastRunAt != nil {
			lastRunAt = row.LastRunAt.ToProto()
		}

		source := &clawv1.SourceData{
			Id:          row.ID,
			Kind:        row.Kind,
			Slug:        row.Slug,
			DisplayName: row.DisplayName,
			Parameter:   row.Parameter,
			Countback:   row.Countback,
			IsDisabled:  bool(row.IsDisabled),
			LastRunAt:   lastRunAt,
			CreatedAt:   row.CreatedAt.ToProto(),
			UpdatedAt:   row.UpdatedAt.ToProto(),
		}
		sources = append(sources, source)

		// Get schedules if requested
		if req.IncludeSchedules {
			schedulesStmt := sqlite.SELECT(table.Schedules.AllColumns).
				FROM(table.Schedules).
				WHERE(table.Schedules.SourceID.EQ(sqlite.Int64(row.ID)))

			var scheduleRows []struct {
				ID        int64 `sql:"primary_key"`
				SourceID  int64
				Schedule  string
				CreatedAt types.UnixMilli
				UpdatedAt types.UnixMilli
			}

			err = schedulesStmt.QueryContext(ctx, s.db, &scheduleRows)
			if err != nil {
				return nil, fmt.Errorf("failed to get schedules for source %d: %w", row.ID, err)
			}

			var schedules []*clawv1.SourceSchedule
			for _, scheduleRow := range scheduleRows {
				schedules = append(schedules, &clawv1.SourceSchedule{
					Id:        scheduleRow.ID,
					SourceId:  scheduleRow.SourceID,
					Schedule:  scheduleRow.Schedule,
					CreatedAt: scheduleRow.CreatedAt.ToProto(),
					UpdatedAt: scheduleRow.UpdatedAt.ToProto(),
				})
			}
			schedulesMap[row.ID] = &clawv1.SourceScheduleList{Schedules: schedules}
		}
	}

	response := &clawv1.ListSourcesResponse{
		Sources:       sources,
		Schedules:     schedulesMap,
		NextPageToken: nextPageToken,
	}

	return response, nil
}

