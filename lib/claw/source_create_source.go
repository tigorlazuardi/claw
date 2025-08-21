package claw

import (
	"context"
	"fmt"

	"github.com/tigorlazuardi/claw/lib/claw/gen/model"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/claw/gen/table"
	"github.com/tigorlazuardi/claw/lib/claw/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CreateSource creates a new source with optional schedules
func (s *Claw) CreateSource(ctx context.Context, req *clawv1.CreateSourceRequest) (*clawv1.CreateSourceResponse, error) {
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
		table.Sources.IsDisabled,
		table.Sources.CreatedAt,
		table.Sources.UpdatedAt,
	).VALUES(
		req.Kind,
		req.Slug,
		req.DisplayName,
		req.Parameter,
		req.Countback,
		types.Bool(req.IsDisabled),
		nowMillis,
		nowMillis,
	).RETURNING(table.Sources.AllColumns)

	var sourceRow model.Sources

	err = sourceStmt.QueryContext(ctx, tx, &sourceRow)
	if err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}

	// Create schedules if provided
	var schedules []*clawv1.SourceSchedule
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

			var scheduleRow model.Schedules

			err = scheduleStmt.QueryContext(ctx, tx, &scheduleRow)
			if err != nil {
				return nil, fmt.Errorf("failed to create schedule: %w", err)
			}

			schedules = append(schedules, &clawv1.SourceSchedule{
				Id:        *scheduleRow.ID,
				SourceId:  scheduleRow.SourceID,
				Schedule:  scheduleRow.Schedule,
				CreatedAt: scheduleRow.CreatedAt.ToProto(),
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

	source := &clawv1.SourceData{
		Id:          *sourceRow.ID,
		Kind:        sourceRow.Kind,
		Slug:        sourceRow.Slug,
		DisplayName: sourceRow.DisplayName,
		Parameter:   sourceRow.Parameter,
		Countback:   int32(sourceRow.Countback),
		IsDisabled:  bool(sourceRow.IsDisabled),
		LastRunAt:   lastRunAt,
		CreatedAt:   sourceRow.CreatedAt.ToProto(),
		UpdatedAt:   sourceRow.UpdatedAt.ToProto(),
	}

	return &clawv1.CreateSourceResponse{
		Source:    source,
		Schedules: schedules,
	}, nil
}
