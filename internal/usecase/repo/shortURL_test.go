package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortURLRepo_Get(t *testing.T) {
	repo := New()

	repo.store["originalURL"] = "shortURL"
	repo.store["shortURL"] = "originalURL"

	type want struct {
		url string
		ok  bool
	}

	cases := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "short url found",
			url:  "originalURL",
			want: want{
				url: "shortURL",
				ok:  true,
			},
		},
		{
			name: "short url not found",
			url:  "exampleURL",
			want: want{
				url: "",
				ok:  false,
			},
		},
		{
			name: "original url found",
			url:  "shortURL",
			want: want{
				url: "originalURL",
				ok:  true,
			},
		},
		{
			name: "original url not found",
			url:  "shortexampleURL",
			want: want{
				url: "",
				ok:  false,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			u, ok := repo.Get(tt.url)
			require.Equal(t, tt.want.ok, ok)
			assert.Equal(t, tt.want.url, u)
		})
	}
}

func TestShortURLRepo_Create(t *testing.T) {
	repo := New()

	type want struct {
		shortURL    string
		originalURL string
	}
	cases := []struct {
		name        string
		originalURL string
		shortURL    string
		want        want
	}{
		{
			name:        "url added to store",
			originalURL: "originalURL",
			shortURL:    "shortURL",
			want: want{
				shortURL:    "shortURL",
				originalURL: "originalURL",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			repo.Create(tt.originalURL, tt.shortURL)

			u, _ := repo.Get(tt.originalURL)
			require.Equal(t, tt.want.shortURL, u)

			u, _ = repo.Get(tt.shortURL)
			require.Equal(t, tt.want.originalURL, u)
		})
	}
}
