package claw

import (
	"context"
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/claw/types"
)

// GetImage retrieves an image by ID with all related data
func (s *Claw) GetImage(ctx context.Context, req *clawv1.GetImageRequest) (*clawv1.GetImageResponse, error) {
	// Get image
	stmt := SELECT(Images.AllColumns).
		FROM(Images).
		WHERE(Images.ID.EQ(Int64(req.Id)))

	var imageRow model.Images
	err := stmt.QueryContext(ctx, s.db, &imageRow)
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	// Get device assignments
	deviceStmt := SELECT(ImageDevices.DeviceID).
		FROM(ImageDevices).
		WHERE(ImageDevices.ImageID.EQ(Int64(req.Id)))

	var deviceRows []model.ImageDevices
	err = deviceStmt.QueryContext(ctx, s.db, &deviceRows)
	if err != nil {
		return nil, fmt.Errorf("failed to get image devices: %w", err)
	}

	var deviceIDs []int64
	for _, device := range deviceRows {
		deviceIDs = append(deviceIDs, device.DeviceID)
	}

	// Get file paths
	pathStmt := SELECT(ImagePaths.Path).
		FROM(ImagePaths).
		WHERE(ImagePaths.ImageID.EQ(Int64(req.Id)))

	var pathRows []model.ImagePaths
	err = pathStmt.QueryContext(ctx, s.db, &pathRows)
	if err != nil {
		return nil, fmt.Errorf("failed to get image paths: %w", err)
	}

	var paths []string
	for _, path := range pathRows {
		paths = append(paths, path.Path)
	}

	// Get tags
	tagStmt := SELECT(ImageTags.Tag).
		FROM(ImageTags).
		WHERE(ImageTags.ImageID.EQ(Int64(req.Id)))

	var tagRows []model.ImageTags
	err = tagStmt.QueryContext(ctx, s.db, &tagRows)
	if err != nil {
		return nil, fmt.Errorf("failed to get image tags: %w", err)
	}

	var tags []string
	for _, tag := range tagRows {
		tags = append(tags, tag.Tag)
	}

	// Convert to protobuf
	image := s.imageModelToProto(imageRow, deviceIDs, paths, tags)

	return &clawv1.GetImageResponse{
		Image: image,
	}, nil
}

// imageModelToProto converts a database image model to protobuf
func (s *Claw) imageModelToProto(imageRow model.Images, deviceIDs []int64, paths []string, tags []string) *clawv1.Image {
	var thumbnailPath *string
	if imageRow.ThumbnailPath != nil {
		thumbnailPath = imageRow.ThumbnailPath
	}

	var postAuthor *string
	if imageRow.PostAuthor != nil {
		postAuthor = imageRow.PostAuthor
	}

	var postAuthorURL *string
	if imageRow.PostAuthorURL != nil {
		postAuthorURL = imageRow.PostAuthorURL
	}

	var postURL *string
	if imageRow.PostURL != nil {
		postURL = imageRow.PostURL
	}

	return &clawv1.Image{
		Id:             *imageRow.ID,
		SourceId:       imageRow.SourceID,
		DeviceIds:      deviceIDs,
		Paths:          paths,
		DownloadUrl:    imageRow.DownloadURL,
		Width:          int32(imageRow.Width),
		Height:         int32(imageRow.Height),
		Filesize:       uint32(imageRow.Filesize),
		ThumbnailPath:  thumbnailPath,
		ImagePath:      imageRow.ImagePath,
		PostAuthor:     postAuthor,
		PostAuthorUrl:  postAuthorURL,
		PostUrl:        postURL,
		IsFavorite:     bool(types.Bool(imageRow.IsFavorite)),
		Tags:           tags,
		CreatedAt:      imageRow.CreatedAt.ToProto(),
		UpdatedAt:      imageRow.UpdatedAt.ToProto(),
	}
}