package reddit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHTTPClient is a mock implementation of the Doer interface
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestIsImgurURL(t *testing.T) {
	reddit := &Reddit{}

	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{"imgur direct link", "https://i.imgur.com/abc123.jpg", true},
		{"imgur gallery", "https://imgur.com/gallery/abc123", true},
		{"imgur with www", "https://www.imgur.com/abc123", true},
		{"imgur uppercase", "https://IMGUR.COM/abc123.jpg", true},
		{"reddit image", "https://i.redd.it/abc123.jpg", false},
		{"other domain", "https://example.com/image.jpg", false},
		{"empty url", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reddit.isImgurURL(tt.url)
			assert.Equal(t, tt.expected, result, "Expected %v for URL %s", tt.expected, tt.url)
		})
	}
}

func TestIsImgurImageValid(t *testing.T) {
	mockClient := &MockHTTPClient{}
	reddit := &Reddit{Client: mockClient}
	ctx := context.Background()

	tests := []struct {
		name           string
		url            string
		statusCode     int
		expectedValid  bool
		expectError    bool
	}{
		{"valid image", "https://i.imgur.com/valid.jpg", 200, true, false},
		{"deleted image", "https://i.imgur.com/deleted.jpg", 404, false, false},
		{"forbidden image", "https://i.imgur.com/forbidden.jpg", 403, false, false},
		{"rate limited", "https://i.imgur.com/ratelimited.jpg", 429, true, false},
		{"server error", "https://i.imgur.com/error.jpg", 500, false, false},
		{"redirect", "https://i.imgur.com/redirect.jpg", 302, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock response
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       http.NoBody,
			}

			if tt.expectError {
				mockClient.On("Do", mock.AnythingOfType("*http.Request")).Return((*http.Response)(nil), assert.AnError).Once()
			} else {
				mockClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()
			}

			result := reddit.isImgurImageValid(ctx, tt.url)
			assert.Equal(t, tt.expectedValid, result, "Expected %v for status code %d", tt.expectedValid, tt.statusCode)

			mockClient.AssertExpectations(t)
		})
	}
}

func TestIsImagePost_WithImgurFiltering(t *testing.T) {
	mockClient := &MockHTTPClient{}
	reddit := &Reddit{Client: mockClient}
	ctx := context.Background()

	tests := []struct {
		name           string
		post           RedditPostData
		imgurStatus    int
		expectedResult bool
		shouldCallAPI  bool
	}{
		{
			name: "valid imgur image",
			post: RedditPostData{
				URL:      "https://i.imgur.com/valid.jpg",
				PostHint: "image",
			},
			imgurStatus:    200,
			expectedResult: true,
			shouldCallAPI:  true,
		},
		{
			name: "deleted imgur image",
			post: RedditPostData{
				URL:      "https://i.imgur.com/deleted.jpg",
				PostHint: "image",
			},
			imgurStatus:    404,
			expectedResult: false,
			shouldCallAPI:  true,
		},
		{
			name: "reddit hosted image",
			post: RedditPostData{
				URL:      "https://i.redd.it/valid.jpg",
				PostHint: "image",
			},
			expectedResult: true,
			shouldCallAPI:  false,
		},
		{
			name: "non-image post",
			post: RedditPostData{
				URL:      "https://example.com/article",
				PostHint: "link",
			},
			expectedResult: false,
			shouldCallAPI:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldCallAPI {
				resp := &http.Response{
					StatusCode: tt.imgurStatus,
					Body:       http.NoBody,
				}
				mockClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()
			}

			result := reddit.isImagePost(ctx, tt.post)
			assert.Equal(t, tt.expectedResult, result, "Expected %v for post %+v", tt.expectedResult, tt.post)

			if tt.shouldCallAPI {
				mockClient.AssertExpectations(t)
			}
		})
	}
}