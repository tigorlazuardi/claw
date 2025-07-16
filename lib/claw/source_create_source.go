package claw

import (
	"context"
	"fmt"

	sourcev1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/source/v1"
	"github.com/tigorlazuardi/claw/lib/claw/gen/table"
	"github.com/tigorlazuardi/claw/lib/claw/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CreateSource creates a new source with optional schedules
func (s *SourceService) CreateSource(ctx context.Context, req *sourcev1.CreateSourceRequest) (*sourcev1.CreateSourceResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	nowMillis := types.UnixMilliNow()

	// Insert source
	sourceStmt := table.Sources.INSERT(
		table.Sources.Kind,
		table.Sources.Slug,
		table.Sources.DisplayName,
		table.Sources.Parameter,
		table.Sources.Countback,
		table.Sources.CreatedAt,
		table.Sources.UpdatedAt,
	).VALUES(
		req.Kind,
		req.Slug,
		req.DisplayName,
		req.Parameter,
		req.Countback,
		nowMillis,
		nowMillis,
	).RETURNING(table.Sources.AllColumns)

	var sourceRow struct {
		ID          int64 `sql:"primary_key"`
		Kind        string
		Slug        string
		DisplayName string
		Parameter   string
		Countback   int32
		LastRunAt   *types.UnixMilli
		CreatedAt   types.UnixMilli
		UpdatedAt   types.UnixMilli
	}

	err = sourceStmt.QueryContext(ctx, tx, &sourceRow)
	if err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}

	// Create schedules if provided
	var schedules []*sourcev1.Schedule
	if len(req.Schedules) > 0 {
		for _, scheduleStr := range req.Schedules {
			scheduleStmt := table.Schedules.INSERT(
				table.Schedules.SourceID,
				table.Schedules.Schedule,
				table.Schedules.CreatedAt,
			).VALUES(
				sourceRow.ID,
				scheduleStr,
				nowMillis,
			).RETURNING(table.Schedules.AllColumns)

			var scheduleRow struct {
				ID        int64 `sql:"primary_key"`
				SourceID  int64
				Schedule  string
				CreatedAt types.UnixMilli
				UpdatedAt types.UnixMilli
			}

			err = scheduleStmt.QueryContext(ctx, tx, &scheduleRow)
			if err != nil {
				return nil, fmt.Errorf("failed to create schedule: %w", err)
			}

			schedules = append(schedules, &sourcev1.Schedule{
				Id:        scheduleRow.ID,
				SourceId:  scheduleRow.SourceID,
				Schedule:  scheduleRow.Schedule,
				CreatedAt: scheduleRow.CreatedAt.ToProto(),
				UpdatedAt: scheduleRow.UpdatedAt.ToProto(),
			})
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
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

	return &sourcev1.CreateSourceResponse{
		Source:    source,
		Schedules: schedules,
	}, nil
}
