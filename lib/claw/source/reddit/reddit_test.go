package reddit

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserPatternRegex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantUser string
		wantType string
		match    bool
	}{
		// User patterns with u/
		{"u/username", "u/testuser", "testuser", "u", true},
		{"u/username with json", "u/testuser.json", "testuser", "u", true},
		{"u/username with trailing slash", "u/testuser/", "testuser", "u", true},
		{"u/username with json and slash", "u/testuser.json/", "testuser", "u", true},

		// User patterns with user/
		{"user/username", "user/testuser", "testuser", "user", true},
		{"user/username with json", "user/testuser.json", "testuser", "user", true},
		{"user/username with trailing slash", "user/testuser/", "testuser", "user", true},
		{"user/username with json and slash", "user/testuser.json/", "testuser", "user", true},

		// Full URL patterns with u/
		{"https://reddit.com/u/username", "https://reddit.com/u/testuser", "testuser", "u", true},
		{"https://reddit.com/u/username.json", "https://reddit.com/u/testuser.json", "testuser", "u", true},
		{"https://reddit.com/u/username/", "https://reddit.com/u/testuser/", "testuser", "u", true},
		{"https://reddit.com/u/username.json/", "https://reddit.com/u/testuser.json/", "testuser", "u", true},

		// Full URL patterns with user/
		{"https://reddit.com/user/username", "https://reddit.com/user/testuser", "testuser", "user", true},
		{"https://reddit.com/user/username.json", "https://reddit.com/user/testuser.json", "testuser", "user", true},
		{"https://reddit.com/user/username/", "https://reddit.com/user/testuser/", "testuser", "user", true},
		{"https://reddit.com/user/username.json/", "https://reddit.com/user/testuser.json/", "testuser", "user", true},

		// www.reddit.com patterns
		{"https://www.reddit.com/u/username", "https://www.reddit.com/u/testuser", "testuser", "u", true},
		{"https://www.reddit.com/user/username", "https://www.reddit.com/user/testuser", "testuser", "user", true},
		{"https://www.reddit.com/u/username.json", "https://www.reddit.com/u/testuser.json", "testuser", "u", true},
		{"https://www.reddit.com/user/username.json", "https://www.reddit.com/user/testuser.json", "testuser", "user", true},

		// http:// patterns
		{"http://reddit.com/u/username", "http://reddit.com/u/testuser", "testuser", "u", true},
		{"http://www.reddit.com/user/username", "http://www.reddit.com/user/testuser", "testuser", "user", true},

		// Complex usernames
		{"username with underscores", "u/test_user", "test_user", "u", true},
		{"username with hyphens", "u/test-user", "test-user", "u", true},
		{"username with numbers", "u/testuser123", "testuser123", "u", true},
		{"username with mixed", "u/test_user-123", "test_user-123", "u", true},

		// Invalid patterns
		{"empty username", "u/", "", "", false},
		{"subreddit pattern", "r/testsubreddit", "", "", false},
		{"invalid characters", "u/test user", "", "", false},
		{"no username prefix", "testuser", "", "", false},
		{"wrong prefix", "p/testuser", "", "", false},
		{"multiple slashes", "u//testuser", "", "", false},
		{"invalid domain", "https://example.com/u/testuser", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := userPattern.FindStringSubmatch(tt.input)
			
			if tt.match {
				require.Len(t, matches, 3, "Expected match with 3 groups for input %q", tt.input)
				assert.Equal(t, tt.wantType, matches[1], "Expected type %q for input %q", tt.wantType, tt.input)
				assert.Equal(t, tt.wantUser, matches[2], "Expected user %q for input %q", tt.wantUser, tt.input)
			} else {
				assert.Empty(t, matches, "Expected no match for input %q", tt.input)
			}
		})
	}
}

func TestSubredditPatternRegex(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantSubreddit string
		match       bool
	}{
		// Basic subreddit patterns
		{"r/subreddit", "r/testsubreddit", "testsubreddit", true},
		{"r/subreddit with json", "r/testsubreddit.json", "testsubreddit", true},
		{"r/subreddit with trailing slash", "r/testsubreddit/", "testsubreddit", true},
		{"r/subreddit with json and slash", "r/testsubreddit.json/", "testsubreddit", true},

		// Full URL patterns
		{"https://reddit.com/r/subreddit", "https://reddit.com/r/testsubreddit", "testsubreddit", true},
		{"https://reddit.com/r/subreddit.json", "https://reddit.com/r/testsubreddit.json", "testsubreddit", true},
		{"https://reddit.com/r/subreddit/", "https://reddit.com/r/testsubreddit/", "testsubreddit", true},
		{"https://reddit.com/r/subreddit.json/", "https://reddit.com/r/testsubreddit.json/", "testsubreddit", true},

		// www.reddit.com patterns
		{"https://www.reddit.com/r/subreddit", "https://www.reddit.com/r/testsubreddit", "testsubreddit", true},
		{"https://www.reddit.com/r/subreddit.json", "https://www.reddit.com/r/testsubreddit.json", "testsubreddit", true},

		// http:// patterns
		{"http://reddit.com/r/subreddit", "http://reddit.com/r/testsubreddit", "testsubreddit", true},
		{"http://www.reddit.com/r/subreddit", "http://www.reddit.com/r/testsubreddit", "testsubreddit", true},

		// Complex subreddit names
		{"subreddit with underscores", "r/test_subreddit", "test_subreddit", true},
		{"subreddit with hyphens", "r/test-subreddit", "test-subreddit", true},
		{"subreddit with numbers", "r/testsubreddit123", "testsubreddit123", true},
		{"subreddit with mixed", "r/test_subreddit-123", "test_subreddit-123", true},

		// Invalid patterns
		{"empty subreddit", "r/", "", false},
		{"user pattern", "u/testuser", "", false},
		{"invalid characters", "r/test subreddit", "", false},
		{"no subreddit prefix", "testsubreddit", "", false},
		{"wrong prefix", "s/testsubreddit", "", false},
		{"multiple slashes", "r//testsubreddit", "", false},
		{"invalid domain", "https://example.com/r/testsubreddit", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := subredditPattern.FindStringSubmatch(tt.input)
			
			if tt.match {
				require.Len(t, matches, 2, "Expected match with 2 groups for input %q", tt.input)
				assert.Equal(t, tt.wantSubreddit, matches[1], "Expected subreddit %q for input %q", tt.wantSubreddit, tt.input)
			} else {
				assert.Empty(t, matches, "Expected no match for input %q", tt.input)
			}
		})
	}
}

func TestRegexPatternCompilation(t *testing.T) {
	t.Run("userPattern compiles without error", func(t *testing.T) {
		pattern := `^(?:https?://(?:www\.)?reddit\.com/)?(u|user)/([a-zA-Z0-9_-]+)(?:\.json)?/?$`
		_, err := regexp.Compile(pattern)
		assert.NoError(t, err, "User pattern should compile without error")
	})

	t.Run("subredditPattern compiles without error", func(t *testing.T) {
		pattern := `^(?:https?://(?:www\.)?reddit\.com/)?r/([a-zA-Z0-9_-]+)(?:\.json)?/?$`
		_, err := regexp.Compile(pattern)
		assert.NoError(t, err, "Subreddit pattern should compile without error")
	})
}

// TestPatternSpecificCases tests edge cases and specific requirements
func TestPatternSpecificCases(t *testing.T) {
	t.Run("user/ prefix gets normalized to u/", func(t *testing.T) {
		matches := userPattern.FindStringSubmatch("user/testuser")
		require.Len(t, matches, 3, "Expected match for user/testuser")
		assert.Equal(t, "user", matches[1], "Expected type 'user'")
		assert.Equal(t, "testuser", matches[2], "Expected username 'testuser'")
	})

	t.Run("both http and https work", func(t *testing.T) {
		httpMatch := userPattern.FindStringSubmatch("http://reddit.com/u/testuser")
		httpsMatch := userPattern.FindStringSubmatch("https://reddit.com/u/testuser")
		
		assert.Len(t, httpMatch, 3, "Expected match for http URL")
		assert.Len(t, httpsMatch, 3, "Expected match for https URL")
	})

	t.Run("case sensitivity", func(t *testing.T) {
		// Test that patterns are case-sensitive for usernames/subreddits but not for domains
		upperUser := userPattern.FindStringSubmatch("u/TestUser")
		lowerUser := userPattern.FindStringSubmatch("u/testuser")
		
		require.Len(t, upperUser, 3, "Upper case username should match")
		require.Len(t, lowerUser, 3, "Lower case username should match")
		
		assert.Equal(t, "TestUser", upperUser[2], "Expected 'TestUser'")
		assert.Equal(t, "testuser", lowerUser[2], "Expected 'testuser'")
	})

	t.Run("special characters in usernames/subreddits", func(t *testing.T) {
		validChars := []string{"test_user", "test-user", "user123", "123user", "a", "A"}
		invalidChars := []string{"test user", "test.user", "test@user", "test#user", "test!user"}
		
		for _, valid := range validChars {
			matches := userPattern.FindStringSubmatch("u/" + valid)
			assert.Len(t, matches, 3, "Expected valid username %q to match", valid)
		}
		
		for _, invalid := range invalidChars {
			matches := userPattern.FindStringSubmatch("u/" + invalid)
			assert.Empty(t, matches, "Expected invalid username %q to not match", invalid)
		}
	})
}