package claw

import (
	"context"
	"fmt"
	"slices"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/claw/types"
)

// ListImages lists images with optional filtering and pagination
func (s *Claw) ListImages(ctx context.Context, req *clawv1.ListImagesRequest) (*clawv1.ListImagesResponse, error) {
	isReversed := req.Pagination != nil && req.Pagination.GetPrevToken() != 0
	cond := Bool(true)
	var from ReadableTable = Images.
		INNER_JOIN(ImagePaths, ImagePaths.ImageID.EQ(Images.ID)).
		INNER_JOIN(ImageTags, ImageTags.ImageID.EQ(Images.ID)).
		INNER_JOIN(ImageDevices, ImageDevices.ImageID.EQ(Images.ID))

	// Search filter
	if search := req.GetSearch(); search != "" {
		searchTerm := "%" + search + "%"
		cond = cond.AND(
			Images.PostAuthor.LIKE(String(searchTerm)).
				OR(Images.PostURL.LIKE(String(searchTerm))).
				OR(Images.DownloadURL.LIKE(String(searchTerm))).
				OR(Images.PostAuthor.LIKE(String(searchTerm))).
				OR(ImageTags.Tag.LIKE(String(searchTerm))),
		)
	}

	// Source filter
	if req.SourceId != nil {
		cond = cond.AND(Images.SourceID.EQ(Int64(*req.SourceId)))
	}

	if req.DeviceId != nil {
		cond = cond.AND(ImageDevices.DeviceID.EQ(Int64(*req.DeviceId)))
	}

	if len(req.Tags) > 0 {
		cond = cond.AND(ImageTags.Tag.IN(jetStringsExpr(req.Tags...)...))
	}

	if req.IsFavorite != nil {
		cond = cond.AND(Images.IsFavorite.EQ(types.NewBoolFromPointer(req.IsFavorite).Integer()))
	}
	limit := int64(50)
	if req.Pagination != nil {
		if token := req.Pagination.GetNextToken(); token != 0 {
			cond = cond.AND(Images.ID.GT(Int64(int64(token))))
		}
		if token := req.Pagination.GetPrevToken(); token != 0 {
			cond = cond.AND(Images.ID.LT(Int64(int64(token))))
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
	// Tiebreaker
	if isReversed {
		sorts = append(sorts, Images.ID.DESC())
	} else {
		sorts = append(sorts, Images.ID.ASC())
	}

	var out []struct {
		model.Images
		ImagePaths   []model.ImagePaths
		ImageDevices []model.ImageDevices
		ImageTags    []model.ImageTags
	}
	err := SELECT(Images.AllColumns).
		FROM(from).
		WHERE(cond).
		ORDER_BY(sorts...).
		LIMIT(limit).
		QueryContext(ctx, s.db, &out)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}
	if len(out) == 0 {
		return &clawv1.ListImagesResponse{
			Images: []*clawv1.Image{},
		}, nil
	}
	if isReversed {
		slices.Reverse(out)
	}
	hasMore := int64(len(out)) >= limit
	var nextPageToken, prevPageToken *uint32
	if hasMore {
		nextPageToken = Ptr(uint32(*out[len(out)-1].ID))
		prevPageToken = Ptr(uint32(*out[0].ID))
	}

	// Convert to []clawv1.Image
	images := make([]*clawv1.Image, len(out))
	for i, row := range out {
		deviceIDs := make([]int64, len(row.ImageDevices))
		for j, device := range row.ImageDevices {
			deviceIDs[j] = device.DeviceID
		}

		paths := make([]string, len(row.ImagePaths))
		for j, path := range row.ImagePaths {
			paths[j] = path.Path
		}

		tags := make([]string, len(row.ImageTags))
		for j, tag := range row.ImageTags {
			tags[j] = tag.Tag
		}

		images[i] = s.imageModelToProto(row.Images, deviceIDs, paths, tags)
	}

	return &clawv1.ListImagesResponse{
		Images: images,
		Pagination: &clawv1.Pagination{
			Size:      Ptr(uint32(len(out))),
			NextToken: nextPageToken,
			PrevToken: prevPageToken,
		},
	}, nil
}
