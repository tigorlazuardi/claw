package claw

import (
	"context"
	"fmt"
	"math"

	. "github.com/go-jet/jet/v2/sqlite"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/claw/source"
)

// getDevicesForJob gets all devices that should receive images for this job
func (c *Claw) getDevicesForJob(ctx context.Context, jobID int64) ([]deviceFilter, error) {
	// Check if job has specific device assignments
	jobImagesQuery := SELECT(JobImages.DeviceID).
		FROM(JobImages).
		WHERE(JobImages.JobID.EQ(Int64(jobID)))

	var jobDeviceIDs []int64
	err := jobImagesQuery.QueryContext(ctx, c.db, &jobDeviceIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get job device assignments: %w", err)
	}

	var devicesResp *clawv1.ListDevicesResponse
	if len(jobDeviceIDs) > 0 {
		// Get specific devices assigned to this job using existing API
		var devices []*clawv1.Device
		for _, deviceID := range jobDeviceIDs {
			deviceResp, err := c.GetDevice(ctx, &clawv1.GetDeviceRequest{Id: deviceID})
			if err != nil {
				continue // Skip devices that can't be found
			}
			devices = append(devices, deviceResp.Device)
		}
		devicesResp = &clawv1.ListDevicesResponse{Devices: devices}
	} else {
		// Get all devices using existing API
		devicesResp, err = c.ListDevices(ctx, &clawv1.ListDevicesRequest{})
		if err != nil {
			return nil, fmt.Errorf("failed to get devices: %w", err)
		}
	}

	var devices []deviceFilter
	for _, device := range devicesResp.Devices {
		devices = append(devices, deviceFilter{
			id:                    device.Id,
			slug:                  device.Slug,
			saveDir:               device.SaveDir,
			width:                 int64(device.Width),
			height:                int64(device.Height),
			aspectRatioDifference: device.AspectRatioDifference,
			imageMinWidth:         int64(device.ImageMinWidth),
			imageMaxWidth:         int64(device.ImageMaxWidth),
			imageMinHeight:        int64(device.ImageMinHeight),
			imageMaxHeight:        int64(device.ImageMaxHeight),
			imageMinFileSize:      int64(device.ImageMinFilesize),
			imageMaxFileSize:      int64(device.ImageMaxFilesize),
			nsfwMode:              int64(device.Nsfw),
		})
	}

	return devices, nil
}

// filterDevicesForImage filters devices that match the given image criteria
func (c *Claw) filterDevicesForImage(img source.Image, devices []deviceFilter) []deviceFilter {
	var matchedDevices []deviceFilter

	for _, device := range devices {
		if c.deviceMatchesImage(img, device) {
			matchedDevices = append(matchedDevices, device)
		}
	}

	return matchedDevices
}

// deviceMatchesImage checks if a device's criteria match the given image
func (c *Claw) deviceMatchesImage(img source.Image, device deviceFilter) bool {
	// Check file size constraints
	if device.imageMinFileSize > 0 && img.Filesize < device.imageMinFileSize {
		return false
	}
	if device.imageMaxFileSize > 0 && img.Filesize > device.imageMaxFileSize {
		return false
	}

	// Check dimension constraints
	if device.imageMinWidth > 0 && img.Width < device.imageMinWidth {
		return false
	}
	if device.imageMaxWidth > 0 && img.Width > device.imageMaxWidth {
		return false
	}
	if device.imageMinHeight > 0 && img.Height < device.imageMinHeight {
		return false
	}
	if device.imageMaxHeight > 0 && img.Height > device.imageMaxHeight {
		return false
	}

	// Check aspect ratio tolerance
	if device.aspectRatioDifference > 0 {
		deviceAspectRatio := float64(device.width) / float64(device.height)
		imageAspectRatio := float64(img.Width) / float64(img.Height)
		aspectDifference := math.Abs(deviceAspectRatio - imageAspectRatio)
		
		if aspectDifference > device.aspectRatioDifference {
			return false
		}
	}

	// NSFW filtering based on device preference
	// Note: This assumes the image has some way to indicate NSFW status
	// For now, we'll accept all images regardless of NSFW status
	// This can be enhanced when NSFW detection is implemented

	return true
}