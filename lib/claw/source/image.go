package source

import "time"

type Image struct {
	// The actual URL to download the image.
	DownloadURL string
	// Width of image in pixels.
	Width int64
	// Height of image in pixels.
	Height int64
	// Filesize of the image in bytes.
	Filesize int64
	// Artist or author of the image, or the uploader's name.
	Author string
	// URL to the author's profile or page.
	AuthorURL string
	// URL to the page where the image is posted so users can find the source.
	Website string
	// Optional thumbnail URL for the image.
	// If empty, the thumbnail will be generated from the image itself.
	ThumbnailURL string
	// When the image was posted or uploaded. Optional.
	PostedAt time.Time
}

type Images []Image
