package server

import (
	"context"

	"connectrpc.com/connect"
	"github.com/tigorlazuardi/claw/lib/claw"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/server/gen/claw/v1/clawv1connect"
)

var _ clawv1connect.DeviceServiceHandler = (*DeviceHandler)(nil)

// DeviceHandler implements the ConnectRPC DeviceService interface
type DeviceHandler struct {
	service *claw.Claw
}

// NewDeviceHandler creates a new DeviceHandler
func NewDeviceHandler(service *claw.Claw) *DeviceHandler {
	return &DeviceHandler{service: service}
}

// CreateDevice handles device creation requests
func (h *DeviceHandler) CreateDevice(ctx context.Context, req *connect.Request[clawv1.CreateDeviceRequest]) (*connect.Response[clawv1.CreateDeviceResponse], error) {
	resp, err := h.service.CreateDevice(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// GetDevice handles device retrieval requests
func (h *DeviceHandler) GetDevice(ctx context.Context, req *connect.Request[clawv1.GetDeviceRequest]) (*connect.Response[clawv1.GetDeviceResponse], error) {
	resp, err := h.service.GetDevice(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// UpdateDevice handles device update requests
func (h *DeviceHandler) UpdateDevice(ctx context.Context, req *connect.Request[clawv1.UpdateDeviceRequest]) (*connect.Response[clawv1.UpdateDeviceResponse], error) {
	resp, err := h.service.UpdateDevice(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// DeleteDevice handles device deletion requests
func (h *DeviceHandler) DeleteDevice(ctx context.Context, req *connect.Request[clawv1.DeleteDeviceRequest]) (*connect.Response[clawv1.DeleteDeviceResponse], error) {
	resp, err := h.service.DeleteDevice(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// ListDevices handles device listing requests
func (h *DeviceHandler) ListDevices(ctx context.Context, req *connect.Request[clawv1.ListDevicesRequest]) (*connect.Response[clawv1.ListDevicesResponse], error) {
	resp, err := h.service.ListDevices(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// ListDropDownDevices handles dropdown device listing requests
func (h *DeviceHandler) ListDropDownDevices(ctx context.Context, req *connect.Request[clawv1.ListDropDownDevicesRequest]) (*connect.Response[clawv1.ListDropDownDevicesResponse], error) {
	resp, err := h.service.ListDropDownDevices(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// SubscribeDevice handles device subscription requests
func (h *DeviceHandler) SubscribeDevice(ctx context.Context, req *connect.Request[clawv1.SubscribeDeviceRequest]) (*connect.Response[clawv1.SubscribeDeviceResponse], error) {
	resp, err := h.service.SubscribeDevice(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// UnsubscribeDevice handles device unsubscription requests
func (h *DeviceHandler) UnsubscribeDevice(ctx context.Context, req *connect.Request[clawv1.UnsubscribeDeviceRequest]) (*connect.Response[clawv1.UnsubscribeDeviceResponse], error) {
	resp, err := h.service.UnsubscribeDevice(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// Ensure DeviceHandler implements the DeviceServiceHandler interface
var _ clawv1connect.DeviceServiceHandler = (*DeviceHandler)(nil)

