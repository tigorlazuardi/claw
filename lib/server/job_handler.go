package server

import (
	"context"

	"connectrpc.com/connect"
	"github.com/tigorlazuardi/claw/lib/claw"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/server/gen/claw/v1/clawv1connect"
)

var _ clawv1connect.JobServiceHandler = (*JobHandler)(nil)

// JobHandler implements the ConnectRPC JobService interface
type JobHandler struct {
	service *claw.Claw
}

// NewJobHandler creates a new JobHandler
func NewJobHandler(service *claw.Claw) *JobHandler {
	return &JobHandler{service: service}
}

// CreateJob handles job creation requests
func (h *JobHandler) CreateJob(ctx context.Context, req *connect.Request[clawv1.CreateJobRequest]) (*connect.Response[clawv1.CreateJobResponse], error) {
	resp, err := h.service.CreateJob(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// GetJob handles job retrieval requests
func (h *JobHandler) GetJob(ctx context.Context, req *connect.Request[clawv1.GetJobRequest]) (*connect.Response[clawv1.GetJobResponse], error) {
	resp, err := h.service.GetJob(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// UpdateJob handles job update requests
func (h *JobHandler) UpdateJob(ctx context.Context, req *connect.Request[clawv1.UpdateJobRequest]) (*connect.Response[clawv1.UpdateJobResponse], error) {
	resp, err := h.service.UpdateJob(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// DeleteJobs handles job deletion requests
func (h *JobHandler) DeleteJobs(ctx context.Context, req *connect.Request[clawv1.DeleteJobsRequest]) (*connect.Response[clawv1.DeleteJobsResponse], error) {
	resp, err := h.service.DeleteJobs(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// ListJobs handles job listing requests
func (h *JobHandler) ListJobs(ctx context.Context, req *connect.Request[clawv1.ListJobsRequest]) (*connect.Response[clawv1.ListJobsResponse], error) {
	resp, err := h.service.ListJobs(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// CancelJob handles job cancellation requests
func (h *JobHandler) CancelJob(ctx context.Context, req *connect.Request[clawv1.CancelJobRequest]) (*connect.Response[clawv1.CancelJobResponse], error) {
	resp, err := h.service.CancelJob(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// RetryJob handles job retry requests
func (h *JobHandler) RetryJob(ctx context.Context, req *connect.Request[clawv1.RetryJobRequest]) (*connect.Response[clawv1.RetryJobResponse], error) {
	resp, err := h.service.RetryJob(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// Ensure JobHandler implements the JobServiceHandler interface
var _ clawv1connect.JobServiceHandler = (*JobHandler)(nil)