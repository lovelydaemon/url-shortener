package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewURL(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		path    string
		want    URL
	}{
		{
			name:    "empty path",
			baseURL: "http://example.com",
			path:    "",
			want:    "http://example.com",
		},
		{
			name:    "path with prefix",
			baseURL: "http://example.com",
			path:    "/test",
			want:    "http://example.com/test",
		},
		{
			name:    "path without prefix",
			baseURL: "http://example.com",
			path:    "test",
			want:    "http://example.com/test",
		},
		{
			name:    "base url with prefix",
			baseURL: "http://example.com",
			path:    "test",
			want:    "http://example.com/test",
		},
		{
			name:    "base url without prefix",
			baseURL: "example.com",
			path:    "test",
			want:    "http://example.com/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := NewURL(tt.baseURL, tt.path)
			assert.Equal(t, tt.want, url)
		})
	}
}

func TestURL_Validate(t *testing.T) {
	tests := []struct {
		name      string
		url       URL
		wantError bool
	}{
		{
			name:      "valid",
			url:       "http://example.com",
			wantError: false,
		},
		{
			name:      "invalid",
			url:       "example.com",
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.url.Validate()
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
