package claw

import (
	"context"
	"fmt"

	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
)

// ValidateSourceParameters validates the given parameter against the specified source backend
// by calling the backend's ValidateTransformParameter method.
//
// This allows the UI to validate parameters before creating or updating sources.
func (claw *Claw) ValidateSourceParameters(ctx context.Context, req *clawv1.ValidateSourceParametersRequest) (*clawv1.ValidateSourceParametersResponse, error) {
	// Find the backend by name
	backend, exists := claw.scheduler.backends[req.SourceName]
	if !exists {
		return nil, fmt.Errorf("source '%s' not found or not registered", req.SourceName)
	}

	// Call the backend's ValidateTransformParameter method
	transformedParam, err := backend.ValidateTransformParameter(ctx, req.Parameter)
	if err != nil {
		return nil, err
	}

	// Parameter is valid, return the transformed parameter
	return &clawv1.ValidateSourceParametersResponse{
		TransformedParameter: transformedParam,
	}, nil
}

