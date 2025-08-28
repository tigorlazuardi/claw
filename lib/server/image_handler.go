package server

import (
	"context"

	"connectrpc.com/connect"
	"github.com/tigorlazuardi/claw/lib/claw"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/server/gen/claw/v1/clawv1connect"
)

// ImageHandler implements the ConnectRPC ImageService interface
type ImageHandler struct {
	service *claw.Claw
}

// NewImageHandler creates a new ImageHandler
func NewImageHandler(service *claw.Claw) *ImageHandler {
	return &ImageHandler{service: service}
}

// GetImage handles image retrieval requests
func (h *ImageHandler) GetImage(ctx context.Context, req *connect.Request[clawv1.GetImageRequest]) (*connect.Response[clawv1.GetImageResponse], error) {
	resp, err := h.service.GetImage(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// ListImages handles image listing requests
func (h *ImageHandler) ListImages(ctx context.Context, req *connect.Request[clawv1.ListImagesRequest]) (*connect.Response[clawv1.ListImagesResponse], error) {
	resp, err := h.service.ListImages(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// UpdateImage handles image update requests
func (h *ImageHandler) UpdateImage(ctx context.Context, req *connect.Request[clawv1.UpdateImageRequest]) (*connect.Response[clawv1.UpdateImageResponse], error) {
	resp, err := h.service.UpdateImage(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// DeleteImages handles image deletion requests
func (h *ImageHandler) DeleteImages(ctx context.Context, req *connect.Request[clawv1.DeleteImagesRequest]) (*connect.Response[clawv1.DeleteImagesResponse], error) {
	resp, err := h.service.DeleteImages(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// MarkFavorite handles image favorite marking requests
func (h *ImageHandler) MarkFavorite(ctx context.Context, req *connect.Request[clawv1.MarkFavoriteRequest]) (*connect.Response[clawv1.MarkFavoriteResponse], error) {
	resp, err := h.service.MarkFavorite(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// AssignTags handles image tag assignment requests
func (h *ImageHandler) AssignTags(ctx context.Context, req *connect.Request[clawv1.AssignTagsRequest]) (*connect.Response[clawv1.AssignTagsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}

// Ensure ImageHandler implements the ImageServiceHandler interface
var _ clawv1connect.ImageServiceHandler = (*ImageHandler)(nil)