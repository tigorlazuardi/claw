package reddit

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/tigorlazuardi/claw/lib/claw/source"
)

const SourceName = "claw.reddit.v1"

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

var _ source.Source = (*Reddit)(nil)

type Reddit struct {
	Client Doer
}

const helpString = /*markdown*/
`This source fetches images from a Reddit user or subreddit.

Supported parameter formats:

- Full URL to a subreddit, e.g. https://reddit.com/r/wallpapers
- Also accept shorthand expression: r/wallpapers.
- u/{user}
- r/{subreddit}
- user/{user} (will be normalized to u/{user})
- user/{user}.json (will be normalized to u/{user})
`

// ParameterHelp returns the help string for the parameter.
// Markdown formatting is supported, but any Javascript will be stripped.
func (re *Reddit) ParameterHelp() string {
	return helpString
}

// ParameterPlaceholder returns the placeholder string for the parameter.
//
// This is usually a very short string to show as a hint for the user.
func (re *Reddit) ParameterPlaceholder() string {
	return `Subreddit name or username, e.g. r/wallpapers or u/spez`
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
	return SourceName
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

// Define regex patterns for validation
var (
	userPattern      = regexp.MustCompile(`^(?:https?://(?:www\.)?reddit\.com/)?(u|user)/([a-zA-Z0-9_-]+)(?:\.json)?/?$`)
	subredditPattern = regexp.MustCompile(`^(?:https?://(?:www\.)?reddit\.com/)?r/([a-zA-Z0-9_-]+)(?:\.json)?/?$`)
)

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
	if param == "" {
		return "", fmt.Errorf("parameter cannot be empty")
	}

	// Check if it matches user pattern
	if matches := userPattern.FindStringSubmatch(param); len(matches) == 3 {
		username := matches[2]
		// Normalize user/<user> to u/<user>
		normalized := "u/" + username

		// Validate against Reddit API and get proper casing
		return re.validateAndNormalizeCasing(ctx, normalized)
	}

	// Check if it matches subreddit pattern
	if matches := subredditPattern.FindStringSubmatch(param); len(matches) == 2 {
		subreddit := matches[1]
		normalized := "r/" + subreddit

		// Validate against Reddit API and get proper casing
		return re.validateAndNormalizeCasing(ctx, normalized)
	}

	// If no patterns match, return error with helpful message
	return "", fmt.Errorf("invalid Reddit parameter format. Supported patterns:\n" +
		"- https://[www.]reddit.com/u/<user>\n" +
		"- https://[www.]reddit.com/u/<user>.json\n" +
		"- https://[www.]reddit.com/user/<user>\n" +
		"- https://[www.]reddit.com/user/<user>.json\n" +
		"- https://[www.]reddit.com/r/<subreddit>\n" +
		"- https://[www.]reddit.com/r/<subreddit>.json\n" +
		"- u/<user>\n" +
		"- u/<user>.json\n" +
		"- r/<subreddit>\n" +
		"- r/<subreddit>.json\n" +
		"- user/<user> (will be normalized to u/<user>)\n" +
		"- user/<user>.json (will be normalized to u/<user>)\n" +
		"\nBracketed means they are optional.",
	)
}

// validateAndNormalizeCasing validates the parameter against Reddit API and normalizes casing
func (re *Reddit) validateAndNormalizeCasing(ctx context.Context, param string) (string, error) {
	// Construct the JSON API URL
	jsonURL := "https://reddit.com/" + param + ".json"

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", jsonURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent header (Reddit requires this)
	req.Header.Set("User-Agent", "claw/1.0")

	// Make the request
	resp, err := re.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to validate parameter against Reddit API: %w", err)
	}
	defer resp.Body.Close()

	// Handle redirects by checking the final URL
	if resp.Request.URL.String() != jsonURL {
		// Parse the redirected URL to get the proper casing
		finalURL := resp.Request.URL.String()
		parsedURL, err := url.Parse(finalURL)
		if err != nil {
			return "", fmt.Errorf("failed to parse redirected URL: %w", err)
		}

		// Extract the path and remove .json suffix if present
		path := parsedURL.Path
		path = strings.TrimSuffix(path, ".json")
		path = strings.TrimSuffix(path, "/")

		// Remove leading slash and return the normalized path
		strings.TrimPrefix(path, "/")

		return path, nil
	}

	// Check if the response is successful
	if resp.StatusCode == 404 {
		if strings.HasPrefix(param, "u/") {
			return "", fmt.Errorf("user '%s' not found on Reddit", param[2:])
		}
		return "", fmt.Errorf("subreddit '%s' not found on Reddit", param[2:])
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Reddit API returned status %d for parameter '%s'", resp.StatusCode, param)
	}

	// If successful and no redirect, return the original normalized parameter
	return param, nil
}
