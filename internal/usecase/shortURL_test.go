package usecase

import (
	"testing"

	"github.com/lovelydaemon/url-shortener/internal/usecase/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortURLUseCase_Get(t *testing.T) {
	usecase := New(repo.New())

	usecase.Create("originalURL", "shortURL")

	type want struct {
		url string
		ok  bool
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "original url found",
			url:  "shortURL",
			want: want{
				url: "originalURL",
				ok:  true,
			},
		},
		{
			name: "short url found",
			url:  "originalURL",
			want: want{
				url: "shortURL",
				ok:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, ok := usecase.Get(tt.url)
			require.Equal(t, tt.want.ok, ok)
			assert.Equal(t, tt.want.url, u)

		})
	}
}
