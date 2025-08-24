package claw

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-jet/jet/v2/sqlite"
	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetSource retrieves a source by ID
func (s *Claw) GetSource(ctx context.Context, req *clawv1.GetSourceRequest) (*clawv1.GetSourceResponse, error) {
	var from ReadableTable = Sources
	if req.GetIncludeSchedules() {
		from = Sources.LEFT_JOIN(Schedules, Sources.ID.EQ(Schedules.SourceID))
	}
	// Get source
	sourceStmt := sqlite.SELECT(Sources.AllColumns).
		FROM(from).
		WHERE(Sources.ID.EQ(sqlite.Int64(req.Id)))

	var sourceRow model.Sources

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

	source := &clawv1.Source{
		Kind:        sourceRow.Kind,
		DisplayName: sourceRow.DisplayName,
		Parameter:   sourceRow.Parameter,
		Countback:   int32(sourceRow.Countback),
		IsDisabled:  bool(sourceRow.IsDisabled),
		LastRunAt:   lastRunAt,
		CreatedAt:   sourceRow.CreatedAt.ToProto(),
		UpdatedAt:   sourceRow.UpdatedAt.ToProto(),
	}

	response := &clawv1.GetSourceResponse{Source: source}

	// Get schedules if requested
	if req.GetIncludeSchedules() {
		schedulesStmt := sqlite.SELECT(Schedules.AllColumns).
			FROM(Schedules).
			WHERE(Schedules.SourceID.EQ(sqlite.Int64(int64(*sourceRow.ID))))

		var scheduleRows []model.Schedules

		err = schedulesStmt.QueryContext(ctx, s.db, &scheduleRows)
		if err != nil {
			return nil, fmt.Errorf("failed to get schedules: %w", err)
		}

		for _, row := range scheduleRows {
			response.Schedules = append(response.Schedules, &clawv1.SourceSchedule{
				Id:        int64(*row.ID),
				Schedule:  row.Schedule,
				CreatedAt: row.CreatedAt.ToProto(),
			})
		}
	}

	return response, nil
}
