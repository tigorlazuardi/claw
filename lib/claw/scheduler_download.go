package claw

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	"github.com/tigorlazuardi/claw/lib/claw/source"
	"github.com/tigorlazuardi/claw/lib/claw/types"
)

// processDownload downloads and processes an image for the given devices
func (scheduler *scheduler) processDownload(ctx context.Context, image source.Image, devices []model.Devices, sourceName string) (err error) {
	// Construct the image file path
	imageDir := filepath.Join(scheduler.config.Download.BaseDir, "images", sourceName)
	imagePath := filepath.Join(imageDir, image.Filename)

	// Check if image already exists
	shouldDownload, err := scheduler.shouldDownloadImage(imagePath)
	if err != nil {
		return fmt.Errorf("failed to check if image should be downloaded: %w", err)
	}

	if shouldDownload {
		// Download image to temporary location first
		tmpPath, err := scheduler.downloadImageToTemp(ctx, image)
		if err != nil {
			return fmt.Errorf("failed to download image: %w", err)
		}
		defer os.Remove(tmpPath) // Clean up temp file

		// Ensure image directory exists
		if err := os.MkdirAll(imageDir, 0o755); err != nil {
			return fmt.Errorf("failed to create image directory: %w", err)
		}

		// Move from temp to final location
		if err := scheduler.moveToFinalLocation(ctx, tmpPath, imagePath); err != nil {
			return fmt.Errorf("failed to move image to final location: %w", err)
		}
	}

	// Find or create image in database
	imageID, err := scheduler.findOrCreateImage(ctx, image, sourceName, imagePath)
	if err != nil {
		return fmt.Errorf("failed to find or create image: %w", err)
	}

	// Process devices and create hardlinks/copies
	for _, device := range devices {
		if err := scheduler.processDeviceAssignment(ctx, image, device, imagePath, sourceName, imageID); err != nil {
			scheduler.logger.ErrorContext(ctx, "failed to process device assignment",
				"device_id", device.ID, "device_name", device.Name, "error", err)
			continue
		}
	}

	return nil
}

// shouldDownloadImage checks if an image should be downloaded
func (scheduler *scheduler) shouldDownloadImage(imagePath string) (bool, error) {
	info, err := os.Stat(imagePath)
	if errors.Is(err, os.ErrNotExist) {
		// Image doesn't exist, should download
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to stat image file: %w", err)
	}

	// Check if file size is suspicious
	if scheduler.config.Download.SanityCheck.Enabled {
		threshold := int64(scheduler.config.Download.SanityCheck.MinImageFilesize)
		if info.Size() < threshold {
			scheduler.logger.InfoContext(context.Background(), "image file size is under threshold, redownloading",
				"path", imagePath, "size", info.Size(), "threshold", threshold)
			return true, nil
		}
	}

	// Image exists and is not suspicious, skip download
	return false, nil
}

// downloadImageToTemp downloads an image to a temporary location
func (scheduler *scheduler) downloadImageToTemp(ctx context.Context, image source.Image) (string, error) {
	// Ensure temp directory exists
	if err := os.MkdirAll(scheduler.config.Download.TmpDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Create temp file
	tmpFile, err := os.CreateTemp(scheduler.config.Download.TmpDir, "claw_download_*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// Download image
	req, err := http.NewRequestWithContext(ctx, "GET", image.DownloadURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create download request: %w", err)
	}

	resp, err := scheduler.httpclient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create stall reader if monitoring is enabled
	var reader io.Reader = resp.Body
	if scheduler.config.Download.StallMonitor.Enabled {
		stallReader := NewStallReader(ctx, resp.Body, scheduler.config.Download.StallMonitor)
		reader = stallReader
	}

	// Copy response body to temp file
	_, err = io.Copy(tmpFile, reader)
	if err != nil {
		return "", fmt.Errorf("failed to copy image data: %w", err)
	}

	return tmpFile.Name(), nil
}

// moveToFinalLocation moves a file from temp location to final location using hardlink or copy
func (scheduler *scheduler) moveToFinalLocation(ctx context.Context, srcPath, dstPath string) error {
	// Try hardlink first
	if err := os.Link(srcPath, dstPath); err == nil {
		scheduler.logger.InfoContext(ctx, "created hardlink for image", "src", srcPath, "dst", dstPath)
		return nil
	}
	scheduler.logger.DebugContext(ctx, "hardlink failed, falling back to copy", "src", srcPath, "dst", dstPath)

	// Hardlink failed, fallback to copy
	return scheduler.copyFile(srcPath, dstPath)
}

// copyFile copies a file from src to dst
func (scheduler *scheduler) copyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// findOrCreateImage finds an existing image or creates a new one in the database
func (scheduler *scheduler) findOrCreateImage(ctx context.Context, image source.Image, sourceName string, imagePath string) (int64, error) {
	// First, try to find existing image by download URL
	var existingImage model.Images
	err := SELECT(Images.AllColumns).
		WHERE(Images.DownloadURL.EQ(String(image.DownloadURL))).
		QueryContext(ctx, scheduler.claw.db, &existingImage)

	if err == nil {
		// Image exists, update the image path and return ID
		relativeImagePath := strings.TrimPrefix(imagePath, scheduler.config.Download.BaseDir+"/")
		_, err = Images.UPDATE(Images.ImagePath, Images.UpdatedAt).
			SET(String(relativeImagePath), types.UnixMilliNow()).
			WHERE(Images.ID.EQ(Int64(*existingImage.ID))).
			ExecContext(ctx, scheduler.claw.db)
		if err != nil {
			return 0, fmt.Errorf("failed to update image path: %w", err)
		}
		return *existingImage.ID, nil
	}

	// Image doesn't exist, create new one
	// First, get source ID for the source name
	var source model.Sources
	err = SELECT(Sources.AllColumns).
		WHERE(Sources.Name.EQ(String(sourceName))).
		QueryContext(ctx, scheduler.claw.db, &source)
	if err != nil {
		return 0, fmt.Errorf("failed to find source: %w", err)
	}

	nowMillis := types.UnixMilliNow()
	relativeImagePath := strings.TrimPrefix(imagePath, scheduler.config.Download.BaseDir+"/")

	// Insert new image
	imageModel := model.Images{
		SourceID:      *source.ID,
		DownloadURL:   image.DownloadURL,
		Width:         image.Width,
		Height:        image.Height,
		Filesize:      image.Filesize,
		ImagePath:     relativeImagePath,
		PostAuthor:    image.Author,
		PostAuthorURL: image.AuthorURL,
		PostURL:       image.Website,
		IsFavorite:    types.Bool(false),
		CreatedAt:     nowMillis,
		UpdatedAt:     nowMillis,
	}

	stmt := Images.INSERT(
		Images.SourceID,
		Images.DownloadURL,
		Images.Width,
		Images.Height,
		Images.Filesize,
		Images.ImagePath,
		Images.PostAuthor,
		Images.PostAuthorURL,
		Images.PostURL,
		Images.IsFavorite,
		Images.CreatedAt,
		Images.UpdatedAt,
	).MODEL(imageModel).
		RETURNING(Images.ID)

	var imageID int64
	err = stmt.QueryContext(ctx, scheduler.claw.db, &imageID)
	if err != nil {
		return 0, fmt.Errorf("failed to create image: %w", err)
	}

	return imageID, nil
}

// processDeviceAssignment handles device assignment and creates hardlinks/copies
func (scheduler *scheduler) processDeviceAssignment(ctx context.Context, image source.Image, device model.Devices, imagePath, sourceName string, imageID int64) error {
	// Generate filename template
	filenameTemplate := device.FilenameTemplate
	if filenameTemplate == "" {
		filenameTemplate = sourceName + "_" + image.Filename
	}

	// Determine target directory
	var targetDir string
	if device.SaveDir == "" {
		targetDir = filepath.Join(scheduler.config.Download.BaseDir, "devices", device.Slug)
	} else {
		targetDir = device.SaveDir
	}

	targetPath := filepath.Join(targetDir, filenameTemplate)

	// Ensure target directory exists
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("failed to create device directory: %w", err)
	}

	// Try hardlink first, fallback to copy
	if err := os.Link(imagePath, targetPath); err != nil {
		if err := scheduler.copyFile(imagePath, targetPath); err != nil {
			return fmt.Errorf("failed to copy image to device location: %w", err)
		}
	}

	// Update ImagePaths table
	relativeDevicePath := strings.TrimPrefix(targetPath, scheduler.config.Download.BaseDir+"/")

	// Insert device path into ImagePaths table
	nowMillis := types.UnixMilliNow()
	_, err := ImagePaths.INSERT(
		ImagePaths.ImageID,
		ImagePaths.Path,
		ImagePaths.CreatedAt,
	).VALUES(
		imageID,
		relativeDevicePath,
		nowMillis,
	).ExecContext(ctx, scheduler.claw.db)
	if err != nil {
		return fmt.Errorf("failed to insert image path: %w", err)
	}

	return nil
}
