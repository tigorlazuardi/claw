package reddit

import (
	"context"
	"net/http"

	"github.com/tigorlazuardi/claw/lib/claw/source"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Reddit struct {
	Client Doer
}

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
//
//	"<namespace>.<name>.<version>"
//
// Namespace can be your name or whatever you want to identify yourself.
//
// Examples:
//
//	"claw.reddit.v1"
//
// Names are tied heavily to parameters. If your parameter schema changes,
// you should also change the name (e.g. bump the version) to avoid
// compatibility issues.
func (re Reddit) Name() string {
	return "claw.reddit.v1"
}

// DisplayName returns the human-readable name for the source.
func (re Reddit) DisplayName() string {
	return "Reddit"
}

// Author returns the author name.
func (re *Reddit) Author() string {
	return "Claw"
}

// AuthorURL returns where the Author can be found or contacted.
func (re *Reddit) AuthorURL() string {
	return "https://github.com/tigorlazuardi/claw"
}

// ValidateTransformParameter validates the parameter for the source and transform the parameter if necessary.
//
// Sources can use this to normalize a parameter and allows more flexible input from the user.
//
// For example, in source "claw.reddit.v1", the Source accepts the following inputs:
//
//   - Full URL to a subreddit, e.g. https://reddit.com/r/wallpapers
//   - Also accept shorthand expression: r/wallpapers.
//   - claw.reddit.v1 also tries to match casing.
//   - If parameter is a user (e.g. u/somebody) -> it will be normalized to user/somebody.
//
// The error message (the .Error() method) must be user friendly and contain all necessary information
// to fix the parameter.
//
// This must return nil error if valid.
func (re *Reddit) ValidateTransformParameter(ctx context.Context, param string) (transformed string, err error) {
	panic("not implemented") // TODO: Implement
}

// Run runs the source to fetch image Metadata based on the given request.
//
// Note that Sources must not download the actual image itself (or only download small part of image to get metadata like dimensions if unavailable in conventional means).
// Sources must only return the metadata and the download URL as [Image] objects.
//
// Claw will handle the downloading after running filters and device assignments.
func (re *Reddit) Run(ctx context.Context, request source.Request) (source.Response, error) {
	panic("not implemented") // TODO: Implement
}
