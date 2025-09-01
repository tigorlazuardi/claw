package claw

import (
	"context"

	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	"github.com/tigorlazuardi/claw/lib/claw/source"
)

// TODO: implement actual download and processing logic
func (scheduler *scheduler) processDownload(ctx context.Context, image source.Image, devices []model.Devices) (err error) {
	return nil
}
