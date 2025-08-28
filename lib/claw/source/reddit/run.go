package reddit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tigorlazuardi/claw/lib/claw/source"
)

// Reddit API response structures
type RedditResponse struct {
	Data RedditData `json:"data"`
}

type RedditData struct {
	After    *string      `json:"after"`
	Children []RedditPost `json:"children"`
}

type RedditPost struct {
	Data RedditPostData `json:"data"`
}

type RedditPostData struct {
	Title      string  `json:"title"`
	URL        string  `json:"url"`
	Author     string  `json:"author"`
	Permalink  string  `json:"permalink"`
	CreatedUTC float64 `json:"created_utc"`
	PostHint   string  `json:"post_hint"`
	Preview    *struct {
		Images []struct {
			Source struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"source"`
		} `json:"images"`
	} `json:"preview"`
}

// Run runs the source to fetch image Metadata based on the given request.
//
// Note that Sources must not download the actual image itself (or only download small part of image to get metadata like dimensions if unavailable in conventional means).
// Sources must only return the metadata and the download URL as [Image] objects.
//
// Claw will handle the downloading after running filters and device assignments.
func (re *Reddit) Run(ctx context.Context, request source.Request) (source.Response, error) {
	// Set default countback if not provided or negative
	countback := request.Countback
	if countback <= 0 {
		countback = 300
	}

	var allImages source.Images
	var next string

	for countback > 0 {
		// Determine limit for this request (max 100)
		limit := min(countback, 100)

		// Fetch posts from Reddit
		posts, nextToken, err := re.fetchRedditPosts(ctx, request.Parameter, limit, next)
		if err != nil {
			return source.Response{}, fmt.Errorf("failed to fetch Reddit posts: %w", err)
		}

		// Convert posts to images
		images := re.filterAndConvertPosts(ctx, posts)
		allImages = append(allImages, images...)

		// Update countback and next token
		countback -= len(posts)
		next = nextToken

		// Stop if no more pages or no next token
		if next == "" {
			break
		}
	}

	return source.Response{Images: allImages}, nil
}

// fetchRedditPosts fetches posts from Reddit API
func (re *Reddit) fetchRedditPosts(ctx context.Context, param string, limit int, after string) ([]RedditPostData, string, error) {
	// Build URL
	baseURL := fmt.Sprintf("https://reddit.com/%s.json", param)
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Add query parameters
	q := u.Query()
	q.Set("limit", strconv.Itoa(limit))
	if after != "" {
		q.Set("after", after)
	}
	u.RawQuery = q.Encode()

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent header (Reddit requires this)
	req.Header.Set("User-Agent", "claw/1.0")

	// Make request
	resp, err := re.Client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("Reddit API returned status %d", resp.StatusCode)
	}

	// Parse response
	var redditResp RedditResponse
	if err := json.NewDecoder(resp.Body).Decode(&redditResp); err != nil {
		return nil, "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract post data
	var posts []RedditPostData
	for _, child := range redditResp.Data.Children {
		posts = append(posts, child.Data)
	}

	// Get next token
	nextToken := ""
	if redditResp.Data.After != nil {
		nextToken = *redditResp.Data.After
	}

	return posts, nextToken, nil
}

// filterAndConvertPosts filters posts for images and converts them to source.Image
func (re *Reddit) filterAndConvertPosts(ctx context.Context, posts []RedditPostData) source.Images {
	var images source.Images

	for _, post := range posts {
		// Check if post is an image
		if !re.isImagePost(ctx, post) {
			continue
		}

		image := re.convertPostToImage(post)
		if image != nil {
			images = append(images, *image)
		}
	}

	return images
}

// isImagePost checks if a Reddit post is an image post
func (re *Reddit) isImagePost(ctx context.Context, post RedditPostData) bool {
	// Check post hint first
	if post.PostHint == "image" {
		// For imgur links, verify the image isn't deleted
		if re.isImgurURL(post.URL) {
			if !re.isImgurImageValid(ctx, post.URL) {
				return false
			}
		}
		return true
	}

	// Check URL for common image extensions
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp"}
	lowerURL := strings.ToLower(post.URL)

	for _, ext := range imageExtensions {
		if strings.HasSuffix(lowerURL, ext) {
			// For imgur links, verify the image isn't deleted
			if re.isImgurURL(post.URL) {
				if !re.isImgurImageValid(ctx, post.URL) {
					return false
				}
			}
			return true
		}
	}

	// Check if it's an imgur direct link
	if strings.Contains(lowerURL, "i.imgur.com") {
		if !re.isImgurImageValid(ctx, post.URL) {
			return false
		}
		return true
	}

	// Check if it's a Reddit hosted image
	if strings.Contains(lowerURL, "i.redd.it") {
		return true
	}

	return false
}

// convertPostToImage converts a Reddit post to a source.Image
func (re *Reddit) convertPostToImage(post RedditPostData) *source.Image {
	image := &source.Image{
		DownloadURL: post.URL,
		Author:      post.Author,
		AuthorURL:   fmt.Sprintf("https://reddit.com/u/%s", post.Author),
		Website:     fmt.Sprintf("https://reddit.com%s", post.Permalink),
		PostedAt:    time.Unix(int64(post.CreatedUTC), 0),
	}

	// Try to get dimensions from preview if available
	if post.Preview != nil && len(post.Preview.Images) > 0 {
		previewImg := post.Preview.Images[0]
		image.Width = int64(previewImg.Source.Width)
		image.Height = int64(previewImg.Source.Height)

		// Use preview URL as thumbnail if available
		if previewImg.Source.URL != "" {
			// Reddit escapes URLs in preview, need to unescape
			unescapedURL := strings.ReplaceAll(previewImg.Source.URL, "&amp;", "&")
			image.ThumbnailURL = unescapedURL
		}
	}

	return image
}

// isImgurURL checks if a URL is from Imgur
func (re *Reddit) isImgurURL(url string) bool {
	lowerURL := strings.ToLower(url)
	return strings.Contains(lowerURL, "imgur.com")
}

// isImgurImageValid checks if an Imgur image is valid (not deleted)
func (re *Reddit) isImgurImageValid(ctx context.Context, imgurURL string) bool {
	// Create a HEAD request to check if the image exists without downloading it
	req, err := http.NewRequestWithContext(ctx, "HEAD", imgurURL, nil)
	if err != nil {
		// If we can't create the request, assume it's invalid
		return false
	}

	// Set User-Agent header
	req.Header.Set("User-Agent", "claw/1.0")

	resp, err := re.Client.Do(req)
	if err != nil {
		// If request fails, assume it's invalid
		return false
	}
	defer resp.Body.Close()

	// Check response status
	switch resp.StatusCode {
	case 200:
		// Image exists and is accessible
		return true
	case 404:
		// Image is deleted or doesn't exist
		return false
	case 403:
		// Image might be blocked or private, consider invalid
		return false
	case 429:
		// Rate limited, assume valid to avoid false negatives
		return true
	default:
		// For other status codes, assume invalid to be safe
		return false
	}
}

