package server

import (
	"context"

	"connectrpc.com/connect"
	"github.com/tigorlazuardi/claw/lib/claw"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/server/gen/claw/v1/clawv1connect"
)

// SourceHandler implements the ConnectRPC SourceService interface
type SourceHandler struct {
	service *claw.Claw
}

// NewSourceHandler creates a new SourceHandler
func NewSourceHandler(service *claw.Claw) *SourceHandler {
	return &SourceHandler{service: service}
}

// CreateSource handles source creation requests
func (h *SourceHandler) CreateSource(ctx context.Context, req *connect.Request[clawv1.CreateSourceRequest]) (*connect.Response[clawv1.CreateSourceResponse], error) {
	resp, err := h.service.CreateSource(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// GetSource handles source retrieval requests
func (h *SourceHandler) GetSource(ctx context.Context, req *connect.Request[clawv1.GetSourceRequest]) (*connect.Response[clawv1.GetSourceResponse], error) {
	resp, err := h.service.GetSource(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// UpdateSource handles source update requests
func (h *SourceHandler) UpdateSource(ctx context.Context, req *connect.Request[clawv1.UpdateSourceRequest]) (*connect.Response[clawv1.UpdateSourceResponse], error) {
	resp, err := h.service.UpdateSource(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// DeleteSource handles source deletion requests
func (h *SourceHandler) DeleteSource(ctx context.Context, req *connect.Request[clawv1.DeleteSourceRequest]) (*connect.Response[clawv1.DeleteSourceResponse], error) {
	resp, err := h.service.DeleteSource(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// ListSources handles source listing requests
func (h *SourceHandler) ListSources(ctx context.Context, req *connect.Request[clawv1.ListSourcesRequest]) (*connect.Response[clawv1.ListSourcesResponse], error) {
	resp, err := h.service.ListSources(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// Ensure SourceHandler implements the SourceServiceHandler interface
var _ clawv1connect.SourceServiceHandler = (*SourceHandler)(nil)

