package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
				path:    "token",
			},
			want: "http://example.com/token",
		},
		{
			name: "baseURL_without_scheme_protocol",
			args: args{
				baseURL: "example.com",
				path:    "token",
			},
			want: "http://example.com/token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := CreateValidURL(tt.args.baseURL, tt.args.path)
			assert.Equal(t, tt.want, url)
		})
	}
}
