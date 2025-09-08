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
	// Names are tied heavily to parameters. If your parameter schema breaking changes,
	// you should also change the name (e.g. bump the version) to avoid
	// compatibility issues.
	Name() string

	// Run runs the source to fetch image Metadata based on the given request.
	//
	// Note that Sources must not download the actual image itself (or only download small part of image to get metadata like dimensions if unavailable in conventional means).
	// Sources must only return the metadata and the download URL as [Image] objects.
	//
	// Claw will handle the downloading after running filters and device assignments.
	Run(ctx context.Context, request Request) (Response, error)

	SourceDescription
	SourceParameter
	SourceCountback
	SourceSchedule
}

type SourceDescription interface {
	// Description returns the human-readable description for the source.
	//
	// Markdown formatting is supported.
	//
	// If empty string is returned, the UI will not show any description.
	Description() string

	// DisplayName returns the human-readable name for the source.
	DisplayName() string

	// Author returns the author name.
	Author() string

	// AuthorURL returns where the Author can be found or contacted.
	AuthorURL() string
}

type SourceParameter interface {
	// RequireParameter returns whether this source requires a parameter.
	//
	// If required, User are not allowed to submit empty parameter.
	RequireParameter() bool
	// ParameterHelp returns the help string for the parameter.
	//
	// It should explain what the parameter is and set expectations on what kind of values are allowed.
	//
	// If the parameter is very complex, show a brief example and forward to external documentation for more details.
	//
	// Markdown formatting is supported.
	//
	// If empty string is returned, the UI will not show any help.
	ParameterHelp() string

	// ParameterPlaceholder returns the placeholder string for the parameter to be shown in the input field.
	//
	// If empty string is returned, the UI will use a very generic placeholder like "Enter parameter...". Which is not very helpful.
	//
	// It's recommended to return a short hint string or example to help the user understand what kind of value is expected
	// without the user have to open the help.
	ParameterPlaceholder() string

	// ValidateTransformParameter validates the parameter for the source and transform the parameter if necessary.
	//
	// Sources can implement this method to normalize a parameter and allows more flexible input from the user.
	//
	// For example, in source "claw.reddit.v1", the Source accepts the following inputs:
	//
	//  - Full URL to a subreddit, e.g. https://reddit.com/r/wallpapers
	//  - Also accept shorthand expression: r/wallpapers.
	//  - claw.reddit.v1 also tries to match casing.
	//  - If parameter is a user (e.g. u/somebody) -> it will be normalized to user/somebody.
	//
	// Then after validation and transformation, the parameter will be normalized to a standard format by claw.reddit.v1.
	// Basically transform e.g. "https://reddit.com/r/wallpapers" and "user/spez" to "r/wallpapers" and "u/spez" respectively.
	//
	// The error message (the .Error() method) must be user friendly and contain all necessary information
	// to fix the parameter.
	//
	// This must return nil error if valid.
	ValidateTransformParameter(ctx context.Context, param string) (transformed string, err error)
}

type SourceCountback interface {
	// DefaultCountback returns the default countback hint for the UI to show
	// to the user.
	//
	// Source implementation should return take notes of API rate limits
	// the source might have, and return a reasonable default value that
	// won't easily hit the rate limit.
	//
	// Note that claw will still pass 0 or negative countback to the Run
	// method if the user specifies that value.
	DefaultCountback() int
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
	// FilenameMaxLength is the maximum allowed length for generated filenames including the extension.
	// If 0 or negative, the Source should use its own default value.
	FilenameMaxLength int
}

type Response struct {
	Images Images
}
