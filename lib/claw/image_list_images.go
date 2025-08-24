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

// ListImages lists images with optional filtering and pagination
func (s *Claw) ListImages(ctx context.Context, req *clawv1.ListImagesRequest) (*clawv1.ListImagesResponse, error) {
	cond := Bool(true)
	var from ReadableTable = Images

	// Search filter
	if req.Search != nil && *req.Search != "" {
		searchTerm := "%" + *req.Search + "%"
		cond.AND(
			Images.PostAuthor.LIKE(String(searchTerm)).
				OR(Images.PostURL.LIKE(String(searchTerm))).
				OR(Images.DownloadURL.LIKE(String(searchTerm))).
				OR(Images.PostAuthor.LIKE(String(searchTerm))),
		)
	}

	// Source filter
	if req.SourceId != nil {
		cond.AND(Images.SourceID.EQ(Int64(*req.SourceId)))
	}

	// Device filter (requires join)
	if req.DeviceId != nil {
		from.INNER_JOIN(ImageDevices, ImageDevices.ImageID.EQ(Images.ID))
		cond.AND(ImageDevices.DeviceID.EQ(Int64(*req.DeviceId)))
	}

	limit := int64(50)
	// Favorite filter
	if req.IsFavorite != nil {
		cond.AND(Images.IsFavorite.EQ(types.NewBoolFromPointer(req.IsFavorite).Integer()))
	}
	if req.Pagination != nil {
		if token := req.Pagination.GetNextToken(); token != 0 {
			cond.AND(Images.ID.GT(Int64(int64(token))))
		}
		if token := req.Pagination.GetPrevToken(); token != 0 {
			cond.AND(Images.ID.LT(Int64(int64(token))))
		}
		if size := req.Pagination.GetSize(); size != 0 {
			limit = Clamp(int64(size), 1, 100)
		}
	}
	sorts := make([]OrderByClause, 0, len(req.Sorts)+1)
	for _, sort := range req.Sorts {
		switch sort.Field {
		case clawv1.ImageField_IMAGE_FIELD_ID:
			sorts = append(sorts, toOrderByClause(Images.ID, sort.Desc))
		case clawv1.ImageField_IMAGE_FIELD_SOURCE_ID:
			sorts = append(sorts, toOrderByClause(Images.SourceID, sort.Desc))
		case clawv1.ImageField_IMAGE_FIELD_POST_URL:
			sorts = append(sorts, toOrderByClause(Images.PostURL, sort.Desc))
		case clawv1.ImageField_IMAGE_FIELD_DOWNLOAD_URL:
			sorts = append(sorts, toOrderByClause(Images.DownloadURL, sort.Desc))
		case clawv1.ImageField_IMAGE_FIELD_WIDTH:
			sorts = append(sorts, toOrderByClause(Images.Width, sort.Desc))
		case clawv1.ImageField_IMAGE_FIELD_FILESIZE:
			sorts = append(sorts, toOrderByClause(Images.Filesize, sort.Desc))
		case clawv1.ImageField_IMAGE_FIELD_HEIGHT:
			sorts = append(sorts, toOrderByClause(Images.Height, sort.Desc))
		case clawv1.ImageField_IMAGE_FIELD_POST_AUTHOR:
			sorts = append(sorts, toOrderByClause(Images.PostAuthor, sort.Desc))
		case clawv1.ImageField_IMAGE_FIELD_IS_FAVORITE:
			sorts = append(sorts, toOrderByClause(Images.IsFavorite, sort.Desc))
		case clawv1.ImageField_IMAGE_FIELD_CREATED_AT:
			sorts = append(sorts, toOrderByClause(Images.CreatedAt, sort.Desc))
		case clawv1.ImageField_IMAGE_FIELD_UPDATED_AT:
			sorts = append(sorts, toOrderByClause(Images.UpdatedAt, sort.Desc))
		default:
			continue
		}
	}
	sorts = append(sorts, Images.ID.ASC()) // Tiebreaker

	// Add pagination
	pageSize := int64(20) // default
	if req.PageSize != nil && *req.PageSize > 0 {
		pageSize = int64(*req.PageSize)
	}

	offset := int64(0)
	if req.PageToken != nil && *req.PageToken > 0 {
		offset = int64(*req.PageToken) * pageSize
	}

	stmt = stmt.LIMIT(pageSize).OFFSET(offset)

	// Execute query
	var imageRows []model.Images
	err := stmt.QueryContext(ctx, s.db, &imageRows)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	// Get related data for all images
	images, err := s.enrichImages(ctx, imageRows)
	if err != nil {
		return nil, fmt.Errorf("failed to enrich images: %w", err)
	}

	// Calculate next page token
	var nextPageToken *uint32
	if len(images) == int(pageSize) {
		nextToken := uint32(offset/pageSize + 1)
		nextPageToken = &nextToken
	}

	return &clawv1.ListImagesResponse{
		Images:        images,
		NextPageToken: nextPageToken,
	}, nil
}

// enrichImages adds device assignments, paths, and tags to images
func (s *Claw) enrichImages(ctx context.Context, imageRows []model.Images) ([]*clawv1.Image, error) {
	if len(imageRows) == 0 {
		return []*clawv1.Image{}, nil
	}

	// Get image IDs
	imageIDs := make([]Expression, len(imageRows))
	imageIDMap := make(map[int64]model.Images)
	for i, imageRow := range imageRows {
		imageIDs[i] = Int64(*imageRow.ID)
		imageIDMap[*imageRow.ID] = imageRow
	}

	// Get all device assignments
	deviceStmt := SELECT(ImageDevices.ImageID, ImageDevices.DeviceID).
		FROM(ImageDevices).
		WHERE(ImageDevices.ImageID.IN(imageIDs...))

	var deviceRows []model.ImageDevices
	err := deviceStmt.QueryContext(ctx, s.db, &deviceRows)
	if err != nil {
		return nil, fmt.Errorf("failed to get image devices: %w", err)
	}

	// Group devices by image ID
	deviceMap := make(map[int64][]int64)
	for _, device := range deviceRows {
		deviceMap[device.ImageID] = append(deviceMap[device.ImageID], device.DeviceID)
	}

	// Get all paths
	pathStmt := SELECT(ImagePaths.ImageID, ImagePaths.Path).
		FROM(ImagePaths).
		WHERE(ImagePaths.ImageID.IN(imageIDs...))

	var pathRows []model.ImagePaths
	err = pathStmt.QueryContext(ctx, s.db, &pathRows)
	if err != nil {
		return nil, fmt.Errorf("failed to get image paths: %w", err)
	}

	// Group paths by image ID
	pathMap := make(map[int64][]string)
	for _, path := range pathRows {
		pathMap[path.ImageID] = append(pathMap[path.ImageID], path.Path)
	}

	// Get all tags
	tagStmt := SELECT(ImageTags.ImageID, ImageTags.Tag).
		FROM(ImageTags).
		WHERE(ImageTags.ImageID.IN(imageIDs...))

	var tagRows []model.ImageTags
	err = tagStmt.QueryContext(ctx, s.db, &tagRows)
	if err != nil {
		return nil, fmt.Errorf("failed to get image tags: %w", err)
	}

	// Group tags by image ID
	tagMap := make(map[int64][]string)
	for _, tag := range tagRows {
		tagMap[tag.ImageID] = append(tagMap[tag.ImageID], tag.Tag)
	}

	// Convert to protobuf
	var images []*clawv1.Image
	for _, imageRow := range imageRows {
		imageID := *imageRow.ID
		image := s.imageModelToProto(
			imageRow,
			deviceMap[imageID],
			pathMap[imageID],
			tagMap[imageID],
		)
		images = append(images, image)
	}

	return images, nil
}
