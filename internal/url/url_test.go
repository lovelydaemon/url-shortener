package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Validate(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		expectedError bool
	}{
		{
			name:          "success",
			url:           "https://example.com",
			expectedError: false,
		},
		{
			name:          "error",
			url:           "example.com",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.url)
			if tt.expectedError {
				assert.Error(t, err, "Expected error")
			} else {
				assert.NoError(t, err, "Expected no error")
			}
		})
	}

}

func Test_CreateValidURL(t *testing.T) {
	type args struct {
		baseURL string
		path    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "baseURL_with_scheme_protocol",
			args: args{
				baseURL: "http://example.com",
				path:    "shortURL",
			},
			want: "http://example.com/shortURL",
		},
		{
			name: "baseURL_without_scheme_protocol",
			args: args{
				baseURL: "example.com",
				path:    "shortURL",
			},
			want: "http://example.com/shortURL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := CreateValidURL(tt.args.baseURL, tt.args.path)
			assert.Equal(t, tt.want, url)
		})
	}
}
