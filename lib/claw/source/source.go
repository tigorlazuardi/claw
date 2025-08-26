package source

import "context"

type Source interface {
	// Name returns the unique kind identifier for the source.
	// Must be unique across all sources.
	//
	// Naming scheme is free, but must be URL and filesystem friendly and should
	// use ASCII characters to minimize compatibility issues.
	//
	// It's very recommended to have some kind of versioning in the name,
	// so that breaking changes can be introduced in the future.
	//
	// For reference, this project uses the following scheme:
	//  "<namespace>.<name>.<version>"
	//
	// Namespace can be your name or whatever you want to identify yourself.
	//
	// Examples:
	//   "claw.reddit.v1"
	//
	// Names are tied heavily to parameters. If your parameter schema changes,
	// you should also change the name (e.g. bump the version) to avoid
	// compatibility issues.
	Name() string

	// DisplayName returns the human-readable name for the source.
	DisplayName() string

	// Author returns the author name.
	Author() string

	// AuthorURL returns where the Author can be found or contacted.
	AuthorURL() string

	// ValidateParameter validates the parameter for the source.
	//
	// The error message (the .Error() method) must be user friendly and contain all necessary information
	// to fix the parameter.
	//
	// This must return nil if the parameter is valid.
	ValidateParameter(param string) error

	// Run runs the source to fetch image Metadata based on the given request.
	//
	// Note that Sources must not download the actual image itself (or only download small part of image to get metadata like dimensions if unavailable in conventional means).
	// Sources must only return the metadata and the download URL as [Image] objects.
	//
	// Claw will handle the downloading after running filters and device assignments.
	Run(ctx context.Context, request Request) (Response, error)
}

type Request struct {
	// Parameter is the source parameter.
	Parameter string
	// Countback is the number of items to lookup for, but not necesarrily the number of images required to be returned.
	// This is meant to limit how far back the Source should look for images.
	//
	// e.g. In claw.reddit.v1 Source, this is the number of posts to look back for, not the number of images to return.
	//
	// If within the Countback range, there are not enough images to return, it's okay to return less images or even zero images.
	//
	// If Countback is 0 or negative, the Source should use it's own default value.
	Countback int
}

type Response struct {
	Images Images
}
