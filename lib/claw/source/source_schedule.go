package source

import (
	"time"

	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
)

type SourceSchedule interface {
	// HaveScheduleConflictCheck returns whether this source has custom schedule conflict check.
	//
	// Usually to avoid schedule too close to each other causing rate limit issues.
	//
	// If true, Claw will call [ScheduleConflictCheck] to check for schedule conflicts.
	HaveScheduleConflictCheck() bool

	// ScheduleConflictCheck returns a warning message if the given schedule has conflicts
	// against this source's requirements.
	//
	// Returning empty string means no conflicts.
	//
	// Return a human friendly message to inform the user about the conflict
	// and the reason why.
	//
	// Markdown formatting is supported in the UI.
	//
	// Received Schedules are already filtered to only contain schedules
	// that belongs to this source (filtered by source name).
	ScheduleConflictCheck(req ScheduleConflictCheckRequest) string
}

type ScheduleConflictCheckRequest struct {
	UserNextRun time.Time
	Schedules   []Schedule
}

type Schedule struct {
	Source    model.Sources
	Schedules []model.Schedules
}
