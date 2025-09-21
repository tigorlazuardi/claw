package claw

import (
	"context"
	"fmt"

	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/claw/types"
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
	sourceStmt := Sources.INSERT(
		Sources.Name,
		Sources.DisplayName,
		Sources.Parameter,
		Sources.Countback,
		Sources.IsDisabled,
		Sources.CreatedAt,
		Sources.UpdatedAt,
	).VALUES(
		req.Name,
		req.DisplayName,
		req.Parameter,
		req.Countback,
		types.Bool(req.IsDisabled),
		nowMillis,
		nowMillis,
	).RETURNING(Sources.AllColumns)

	var sourceRow model.Sources

	err = sourceStmt.QueryContext(ctx, tx, &sourceRow)
	if err != nil {
		return nil, fmt.Errorf("failed to create source: %w", err)
	}

	// Create schedules if provided
	var schedules []*clawv1.SourceSchedule
	if len(req.Schedules) > 0 {
		var entries []model.Schedules
		for _, scheduleStr := range req.Schedules {
			entries = append(entries, model.Schedules{
				SourceID:  *sourceRow.ID,
				Schedule:  scheduleStr,
				CreatedAt: nowMillis,
			})
		}
		insert := Schedules.
			INSERT(Schedules.SourceID, Schedules.Schedule, Schedules.CreatedAt).
			MODELS(entries).
			RETURNING(Schedules.AllColumns)

		var out []model.Schedules
		err := insert.QueryContext(ctx, tx, &out)
		if err != nil {
			return nil, fmt.Errorf("failed to create schedules: %w", err)
		}

		for _, scheduleRow := range out {
			schedules = append(schedules, &clawv1.SourceSchedule{
				Schedule:  scheduleRow.Schedule,
				CreatedAt: scheduleRow.CreatedAt.ToProto(),
			})
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	source := &clawv1.Source{
		Name:        sourceRow.Name,
		DisplayName: sourceRow.DisplayName,
		Parameter:   sourceRow.Parameter,
		Countback:   int32(sourceRow.Countback),
		IsDisabled:  bool(sourceRow.IsDisabled),
		LastRunAt:   sourceRow.LastRunAt.ToProto(),
		CreatedAt:   sourceRow.CreatedAt.ToProto(),
		UpdatedAt:   sourceRow.UpdatedAt.ToProto(),
	}

	return &clawv1.CreateSourceResponse{
		Source:    source,
		Schedules: schedules,
	}, nil
}
