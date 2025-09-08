package claw

import (
	"context"

	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
)

// ListAvailableSources returns all registered source types that can be added by users.
// This includes all sources registered in the scheduler.backends map.
func (claw *Claw) ListAvailableSources(ctx context.Context, req *clawv1.ListAvailableSourcesRequest) (*clawv1.ListAvailableSourcesResponse, error) {
	var availableSources []*clawv1.AvailableSource

	// Iterate through all registered backends
	for _, backend := range claw.scheduler.backends {
		availableSource := &clawv1.AvailableSource{
			Name:                      backend.Name(),
			DisplayName:               backend.DisplayName(),
			Author:                    backend.Author(),
			AuthorUrl:                 backend.AuthorURL(),
			ParameterHelp:             backend.ParameterHelp(),
			ParameterPlaceholder:      backend.ParameterPlaceholder(),
			RequireParameter:          backend.RequireParameter(),
			Description:               backend.Description(),
			DefaultCountback:          int32(backend.DefaultCountback()),
			HaveScheduleConflictCheck: backend.HaveScheduleConflictCheck(),
		}
		availableSources = append(availableSources, availableSource)
	}

	return &clawv1.ListAvailableSourcesResponse{
		Sources: availableSources,
	}, nil
}
