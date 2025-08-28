package server

import (
	"context"

	"connectrpc.com/connect"
	"github.com/tigorlazuardi/claw/lib/claw"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/server/gen/claw/v1/clawv1connect"
)

// TagHandler implements the ConnectRPC TagService interface
// Currently returns unimplemented errors as tag functionality is not yet implemented in the claw service
type TagHandler struct {
	service *claw.Claw
}

// NewTagHandler creates a new TagHandler
func NewTagHandler(service *claw.Claw) *TagHandler {
	return &TagHandler{service: service}
}

// CreateTag handles tag creation requests
func (h *TagHandler) CreateTag(ctx context.Context, req *connect.Request[clawv1.CreateTagRequest]) (*connect.Response[clawv1.CreateTagResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}

// GetTag handles tag retrieval requests
func (h *TagHandler) GetTag(ctx context.Context, req *connect.Request[clawv1.GetTagRequest]) (*connect.Response[clawv1.GetTagResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}

// UpdateTag handles tag update requests
func (h *TagHandler) UpdateTag(ctx context.Context, req *connect.Request[clawv1.UpdateTagRequest]) (*connect.Response[clawv1.UpdateTagResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}

// DeleteTags handles tag deletion requests
func (h *TagHandler) DeleteTags(ctx context.Context, req *connect.Request[clawv1.DeleteTagsRequest]) (*connect.Response[clawv1.DeleteTagsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}

// ListTags handles tag listing requests
func (h *TagHandler) ListTags(ctx context.Context, req *connect.Request[clawv1.ListTagsRequest]) (*connect.Response[clawv1.ListTagsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}

// GetPopularTags handles popular tag retrieval requests
func (h *TagHandler) GetPopularTags(ctx context.Context, req *connect.Request[clawv1.GetPopularTagsRequest]) (*connect.Response[clawv1.GetPopularTagsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}

// Ensure TagHandler implements the TagServiceHandler interface
var _ clawv1connect.TagServiceHandler = (*TagHandler)(nil)