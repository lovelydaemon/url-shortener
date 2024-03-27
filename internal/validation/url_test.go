package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidUrl(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid url",
			url:     "https://google.com",
			wantErr: false,
		},
		{
			name:    "invalid url",
			url:     "invalid.com",
			wantErr: true,
		},
		{
			name:    "empty url",
			url:     "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidUrl(tt.url)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
