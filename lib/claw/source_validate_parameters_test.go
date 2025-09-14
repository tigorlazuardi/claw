package claw

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
	"github.com/tigorlazuardi/claw/lib/claw/source"
	"github.com/tigorlazuardi/claw/lib/claw/source/reddit"
)

func TestValidateSourceParameters(t *testing.T) {
	// Create a test Claw instance
	claw := &Claw{
		scheduler: &scheduler{
			backends: map[string]source.Source{
				"claw.reddit.v1": &reddit.Reddit{
					Client: http.DefaultClient,
				},
			},
		},
	}

	tests := []struct {
		name           string
		sourceName     string
		parameter      string
		expectedValid  bool
		expectedError  string
	}{
		{
			name:          "valid reddit subreddit",
			sourceName:    "claw.reddit.v1",
			parameter:     "r/wallpapers",
			expectedValid: true,
		},
		{
			name:          "valid reddit user format",
			sourceName:    "claw.reddit.v1", 
			parameter:     "u/test",
			expectedValid: true,
		},
		{
			name:          "invalid reddit parameter",
			sourceName:    "claw.reddit.v1",
			parameter:     "invalid-format",
			expectedValid: false,
			expectedError: "invalid Reddit parameter format",
		},
		{
			name:          "non-existent source",
			sourceName:    "non.existent.source",
			parameter:     "test",
			expectedValid: false,
			expectedError: "source 'non.existent.source' not found or not registered",
		},
		{
			name:          "empty parameter",
			sourceName:    "claw.reddit.v1",
			parameter:     "",
			expectedValid: false,
			expectedError: "parameter cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &clawv1.ValidateSourceParametersRequest{
				SourceName: tt.sourceName,
				Parameter:  tt.parameter,
			}

			resp, err := claw.ValidateSourceParameters(context.Background(), req)
			require.NoError(t, err)
			require.NotNil(t, resp)

			assert.Equal(t, tt.expectedValid, resp.Valid)
			
			if tt.expectedValid {
				assert.NotEmpty(t, resp.TransformedParameter)
				assert.Empty(t, resp.ErrorMessage)
			} else {
				assert.Empty(t, resp.TransformedParameter)
				assert.Contains(t, resp.ErrorMessage, tt.expectedError)
			}
		})
	}
}