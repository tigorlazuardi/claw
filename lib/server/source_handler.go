package server

import (
	"context"

	"connectrpc.com/connect"
	"github.com/adhocore/gronx"
	"github.com/tigorlazuardi/claw/lib/claw"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/server/gen/claw/v1/clawv1connect"
	"google.golang.org/protobuf/types/known/timestamppb"
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

// ListAvailableSources handles listing available source types
func (h *SourceHandler) ListAvailableSources(ctx context.Context, req *connect.Request[clawv1.ListAvailableSourcesRequest]) (*connect.Response[clawv1.ListAvailableSourcesResponse], error) {
	resp, err := h.service.ListAvailableSources(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// GetCronNextTime handles cron expression next time calculation
func (h *SourceHandler) GetCronNextTime(ctx context.Context, req *connect.Request[clawv1.GetCronNextTimeRequest]) (*connect.Response[clawv1.GetCronNextTimeResponse], error) {
	gron := gronx.New()

	// Validate the cron expression
	if !gron.IsValid(req.Msg.CronExpression) {
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}

	// Get the next run time
	nextTime, err := gronx.NextTick(req.Msg.CronExpression, true)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := &clawv1.GetCronNextTimeResponse{
		NextTime: timestamppb.New(nextTime),
	}

	return connect.NewResponse(resp), nil
}

// Ensure SourceHandler implements the SourceServiceHandler interface
var _ clawv1connect.SourceServiceHandler = (*SourceHandler)(nil)
