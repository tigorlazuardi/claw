package claw

import (
	"context"
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/tigorlazuardi/claw/lib/claw/gen/jet/model"
	. "github.com/tigorlazuardi/claw/lib/claw/gen/jet/table"
	"github.com/tigorlazuardi/claw/lib/claw/source"
	"github.com/tigorlazuardi/claw/lib/claw/types"
)

// downloadWorker processes download tasks from the download queue
func (c *Claw) downloadWorker(ctx context.Context, workerID int) {
	var logger *slog.Logger
	if c.logger != nil {
		logger = c.logger.With("worker_id", workerID, "worker_type", "download")
		logger.Debug("Starting download worker")
		defer logger.Debug("Download worker stopped")
	}

	for {
		select {
		case <-ctx.Done():
			return
		case task := <-c.downloadQueue:
			if logger != nil {
				logger = logger.With("job_id", task.jobID, "download_url", task.image.DownloadURL)
				logger.Debug("Processing download task")
			}

			if err := c.processDownload(ctx, task); err != nil && logger != nil {
				logger.Error("Failed to process download", "error", err)
			}
		}
	}
}

// processDownload handles the download and file management for a single image
func (c *Claw) processDownload(ctx context.Context, task downloadTask) error {
	// Check if image already exists in database
	existingImage, err := c.getExistingImage(ctx, task.image.DownloadURL)
	if err != nil {
		return fmt.Errorf("failed to check existing image: %w", err)
	}

	var imageID int64
	var filename string

	if existingImage != nil {
		// Image exists, check if file exists on filesystem
		imageID = *existingImage.ID
		filename = c.generateFilename(task.image.DownloadURL)
		imagePath := filepath.Join(c.schedulerConfig.BaseDir, "images", task.sourceName, filename)

		if _, err := os.Stat(imagePath); err == nil {
			if c.logger != nil {
				c.logger.Debug("Image file already exists, skipping download")
			}
			return c.handleExistingImage(ctx, imageID, task, imagePath)
		}

		if c.logger != nil {
			c.logger.Debug("Image record exists but file missing, re-downloading")
		}
	}

	// Generate unique filename
	filename = c.generateFilename(task.image.DownloadURL)

	// Download to temporary location
	tmpPath := filepath.Join(c.schedulerConfig.TmpDir, "claw", filename)
	if err := os.MkdirAll(filepath.Dir(tmpPath), 0o755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	if err := c.downloadFile(ctx, task.image.DownloadURL, tmpPath); err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer os.Remove(tmpPath) // Clean up temp file

	// Ensure base images directory exists
	imagesDir := filepath.Join(c.schedulerConfig.BaseDir, "images", task.sourceName)
	if err := os.MkdirAll(imagesDir, 0o755); err != nil {
		return fmt.Errorf("failed to create images directory: %w", err)
	}

	// Copy to final location
	finalPath := filepath.Join(imagesDir, filename)
	if err := c.copyFile(tmpPath, finalPath); err != nil {
		return fmt.Errorf("failed to copy file to final location: %w", err)
	}

	// Save or update image record
	if existingImage == nil {
		imageID, err = c.saveImageRecord(ctx, task.image, finalPath, task.sourceID)
		if err != nil {
			return fmt.Errorf("failed to save image record: %w", err)
		}
	}

	// Handle device assignments and file copying
	return c.handleDeviceAssignments(ctx, imageID, task, finalPath)
}

// getExistingImage checks if an image with the given download URL already exists
func (c *Claw) getExistingImage(ctx context.Context, downloadURL string) (*model.Images, error) {
	query := SELECT(Images.AllColumns).
		FROM(Images).
		WHERE(Images.DownloadURL.EQ(String(downloadURL)))

	var imageRow model.Images
	err := query.QueryContext(ctx, c.db, &imageRow)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &imageRow, nil
}

// generateFilename generates a unique filename based on the download URL
func (c *Claw) generateFilename(downloadURL string) string {
	// Extract file extension from URL
	ext := filepath.Ext(downloadURL)
	if ext == "" {
		ext = ".jpg" // Default extension
	}

	// Generate hash-based filename
	hash := md5.Sum([]byte(downloadURL))
	return fmt.Sprintf("%x%s", hash, ext)
}

// downloadFile downloads a file from URL to the specified path
func (c *Claw) downloadFile(ctx context.Context, url, path string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}

// copyFile copies a file from src to dst
func (c *Claw) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// hardLink attempts to create a hard link, falls back to copy if not possible
func (c *Claw) hardLink(src, dst string) error {
	// Try hard link first
	if err := os.Link(src, dst); err == nil {
		return nil
	}

	// Fall back to copy
	return c.copyFile(src, dst)
}

// saveImageRecord saves a new image record to the database
func (c *Claw) saveImageRecord(ctx context.Context, img source.Image, path string, sourceID int64) (int64, error) {
	nowMillis := types.UnixMilliNow()

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
		Images.CreatedAt,
		Images.UpdatedAt,
	).MODEL(model.Images{
		SourceID:      sourceID,
		DownloadURL:   img.DownloadURL,
		Width:         img.Width,
		Height:        img.Height,
		Filesize:      img.Filesize,
		ImagePath:     path,
		PostAuthor:    img.Author,
		PostAuthorURL: img.AuthorURL,
		PostURL:       img.Website,
		CreatedAt:     nowMillis,
		UpdatedAt:     nowMillis,
	}).RETURNING(Images.ID)

	var imageID int64
	err := stmt.QueryContext(ctx, c.db, &imageID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert image: %w", err)
	}

	// Insert image path
	pathStmt := ImagePaths.INSERT(
		ImagePaths.ImageID,
		ImagePaths.Path,
		ImagePaths.CreatedAt,
	).MODEL(model.ImagePaths{
		ImageID:   imageID,
		Path:      path,
		CreatedAt: nowMillis,
	})

	_, err = pathStmt.ExecContext(ctx, c.db)
	if err != nil {
		return 0, fmt.Errorf("failed to insert image path: %w", err)
	}

	return imageID, nil
}

// handleExistingImage handles copying/linking for existing images
func (c *Claw) handleExistingImage(ctx context.Context, imageID int64, task downloadTask, imagePath string) error {
	return c.handleDeviceAssignments(ctx, imageID, task, imagePath)
}

// handleDeviceAssignments copies/links images to device directories and updates assignments
func (c *Claw) handleDeviceAssignments(ctx context.Context, imageID int64, task downloadTask, sourcePath string) error {
	nowMillis := types.UnixMilliNow()
	filename := filepath.Base(sourcePath)

	for _, device := range task.devices {
		// Determine save directory
		saveDir := device.saveDir
		if saveDir == "" {
			saveDir = filepath.Join(c.schedulerConfig.BaseDir, device.slug)
		}

		// Ensure device directory exists
		if err := os.MkdirAll(saveDir, 0o755); err != nil {
			return fmt.Errorf("failed to create device directory: %w", err)
		}

		// Generate device-specific filename
		deviceFilename := fmt.Sprintf("%s_%s", task.sourceName, filename)
		devicePath := filepath.Join(saveDir, deviceFilename)

		// Copy/hardlink to device directory
		if err := c.hardLink(sourcePath, devicePath); err != nil {
			return fmt.Errorf("failed to copy image to device directory: %w", err)
		}

		// Insert/update image device assignment
		imageDeviceStmt := ImageDevices.INSERT(
			ImageDevices.ImageID,
			ImageDevices.DeviceID,
			ImageDevices.CreatedAt,
		).MODEL(model.ImageDevices{
			ImageID:   imageID,
			DeviceID:  device.id,
			CreatedAt: nowMillis,
		}).ON_CONFLICT(ImageDevices.ImageID, ImageDevices.DeviceID).DO_NOTHING()

		_, err := imageDeviceStmt.ExecContext(ctx, c.db)
		if err != nil {
			return fmt.Errorf("failed to insert image device assignment: %w", err)
		}

		// Insert image path for device
		pathStmt := ImagePaths.INSERT(
			ImagePaths.ImageID,
			ImagePaths.Path,
			ImagePaths.CreatedAt,
		).MODEL(model.ImagePaths{
			ImageID:   imageID,
			Path:      devicePath,
			CreatedAt: nowMillis,
		})

		_, err = pathStmt.ExecContext(ctx, c.db)
		if err != nil {
			return fmt.Errorf("failed to insert device image path: %w", err)
		}
	}

	return nil
}

