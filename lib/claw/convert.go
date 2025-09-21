package claw

import (
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/claw/types"
)

// imageModelToProto converts a database image model to protobuf
func imageModelToProto(imageRow model.Images) *clawv1.Image {
	return &clawv1.Image{
		Id:            *imageRow.ID,
		SourceId:      imageRow.SourceID,
		DownloadUrl:   imageRow.DownloadURL,
		Width:         int32(imageRow.Width),
		Height:        int32(imageRow.Height),
		Filesize:      uint32(imageRow.Filesize),
		ThumbnailPath: &imageRow.ThumbnailPath,
		ImagePath:     imageRow.ImagePath,
		Title:         &imageRow.Title,
		PostAuthor:    &imageRow.PostAuthor,
		PostAuthorUrl: &imageRow.PostAuthorURL,
		PostUrl:       &imageRow.PostURL,
		IsFavorite:    bool(types.Bool(imageRow.IsFavorite)),
		CreatedAt:     imageRow.CreatedAt.ToProto(),
		UpdatedAt:     imageRow.UpdatedAt.ToProto(),
	}
}

func tagModelToProto(tag model.Tags) *clawv1.Tag {
	return &clawv1.Tag{
		Id:        *tag.ID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt.ToProto(),
	}
}

func tagModelsToProto(tags []model.Tags) []*clawv1.Tag {
	out := make([]*clawv1.Tag, 0, len(tags))
	for _, tag := range tags {
		out = append(out, tagModelToProto(tag))
	}
	return out
}

func deviceModelToProto(device model.Devices) *clawv1.Device {
	return &clawv1.Device{
		Id:                    *device.ID,
		Name:                  device.Name,
		CreatedAt:             device.CreatedAt.ToProto(),
		Height:                int32(device.Height),
		Width:                 int32(device.Width),
		AspectRatioDifference: device.AspectRatioDifference,
		FilenameTemplate:      new(string),
		ImageMinHeight:        uint32(device.ImageMinHeight),
		ImageMinWidth:         uint32(device.ImageMinWidth),
		ImageMaxHeight:        uint32(device.ImageMaxHeight),
		ImageMaxWidth:         uint32(device.ImageMaxWidth),
		ImageMinFilesize:      uint32(device.ImageMinFileSize),
		ImageMaxFilesize:      uint32(device.ImageMaxFileSize),
		Nsfw:                  clawv1.NSFWMode(device.NsfwMode),
		UpdatedAt:             device.UpdatedAt.ToProto(),
		LastActiveAt:          device.LastActiveAt.ToProto(),
	}
}

func sourceModelToProto(source model.Sources) *clawv1.Source {
	return &clawv1.Source{
		Id:          *source.ID,
		Name:        source.Name,
		CreatedAt:   source.CreatedAt.ToProto(),
		DisplayName: source.Name,
		Parameter:   "",
		Countback:   0,
		IsDisabled:  false,
		LastRunAt:   source.LastRunAt.ToProto(),
		UpdatedAt:   source.UpdatedAt.ToProto(),
	}
}
