package claw

import (
	"context"
	"database/sql"
	"fmt"

	sourcev1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/source/v1"
	"github.com/tigorlazuardi/claw/lib/claw/gen/table"
	"github.com/tigorlazuardi/claw/lib/claw/types"
	"github.com/go-jet/jet/v2/sqlite"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetSource retrieves a source by ID
func (s *SourceService) GetSource(ctx context.Context, req *sourcev1.GetSourceRequest) (*sourcev1.GetSourceResponse, error) {
	// Get source
	sourceStmt := sqlite.SELECT(table.Sources.AllColumns).
		FROM(table.Sources).
		WHERE(table.Sources.ID.EQ(sqlite.Int64(req.Id)))

	var sourceRow struct {
		ID          int64             `sql:"primary_key"`
		Kind        string            
		Slug        string            
		DisplayName string            
		Parameter   string            
		Countback   int32             
		LastRunAt   *types.UnixMilli  
		CreatedAt   types.UnixMilli   
		UpdatedAt   types.UnixMilli   
	}

	err := sourceStmt.QueryContext(ctx, s.db, &sourceRow)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("source not found")
		}
		return nil, fmt.Errorf("failed to get source: %w", err)
	}

	// Convert to protobuf
	var lastRunAt *timestamppb.Timestamp
	if sourceRow.LastRunAt != nil {
		lastRunAt = sourceRow.LastRunAt.ToProto()
	}

	source := &sourcev1.Source{
		Id:          sourceRow.ID,
		Kind:        sourceRow.Kind,
		Slug:        sourceRow.Slug,
		DisplayName: sourceRow.DisplayName,
		Parameter:   sourceRow.Parameter,
		Countback:   sourceRow.Countback,
		LastRunAt:   lastRunAt,
		CreatedAt:   sourceRow.CreatedAt.ToProto(),
		UpdatedAt:   sourceRow.UpdatedAt.ToProto(),
	}

	response := &sourcev1.GetSourceResponse{Source: source}

	// Get schedules if requested
	if req.IncludeSchedules {
		schedulesStmt := sqlite.SELECT(table.Schedules.AllColumns).
			FROM(table.Schedules).
			WHERE(table.Schedules.SourceID.EQ(sqlite.Int64(req.Id)))

		var scheduleRows []struct {
			ID        int64            `sql:"primary_key"`
			SourceID  int64            
			Schedule  string           
			CreatedAt types.UnixMilli  
			UpdatedAt types.UnixMilli  
		}

		err = schedulesStmt.QueryContext(ctx, s.db, &scheduleRows)
		if err != nil {
			return nil, fmt.Errorf("failed to get schedules: %w", err)
		}

		for _, row := range scheduleRows {
			response.Schedules = append(response.Schedules, &sourcev1.Schedule{
				Id:        row.ID,
				SourceId:  row.SourceID,
				Schedule:  row.Schedule,
				CreatedAt: row.CreatedAt.ToProto(),
				UpdatedAt: row.UpdatedAt.ToProto(),
			})
		}
	}

	return response, nil
}