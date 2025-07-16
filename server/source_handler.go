package server

import (
	"context"

	"connectrpc.com/connect"
	"github.com/tigorlazuardi/claw/lib/claw"
	sourcev1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/source/v1"
	"github.com/tigorlazuardi/claw/server/gen/source/v1/sourcev1connect"
)

// SourceHandler implements the ConnectRPC SourceService interface
type SourceHandler struct {
	service *claw.SourceService
}

// NewSourceHandler creates a new SourceHandler
func NewSourceHandler(service *claw.SourceService) *SourceHandler {
	return &SourceHandler{service: service}
}

// CreateSource handles source creation requests
func (h *SourceHandler) CreateSource(ctx context.Context, req *connect.Request[sourcev1.CreateSourceRequest]) (*connect.Response[sourcev1.CreateSourceResponse], error) {
	resp, err := h.service.CreateSource(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// GetSource handles source retrieval requests
func (h *SourceHandler) GetSource(ctx context.Context, req *connect.Request[sourcev1.GetSourceRequest]) (*connect.Response[sourcev1.GetSourceResponse], error) {
	resp, err := h.service.GetSource(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// UpdateSource handles source update requests
func (h *SourceHandler) UpdateSource(ctx context.Context, req *connect.Request[sourcev1.UpdateSourceRequest]) (*connect.Response[sourcev1.UpdateSourceResponse], error) {
	resp, err := h.service.UpdateSource(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// DeleteSource handles source deletion requests
func (h *SourceHandler) DeleteSource(ctx context.Context, req *connect.Request[sourcev1.DeleteSourceRequest]) (*connect.Response[sourcev1.DeleteSourceResponse], error) {
	resp, err := h.service.DeleteSource(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// ListSources handles source listing requests
func (h *SourceHandler) ListSources(ctx context.Context, req *connect.Request[sourcev1.ListSourcesRequest]) (*connect.Response[sourcev1.ListSourcesResponse], error) {
	resp, err := h.service.ListSources(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// Ensure SourceHandler implements the SourceServiceHandler interface
var _ sourcev1connect.SourceServiceHandler = (*SourceHandler)(nil)