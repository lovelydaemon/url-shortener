package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
	"github.com/lovelydaemon/url-shortener/internal/usecase/repo"
	"github.com/stretchr/testify/assert"
)

func Test_ShortenRoutes_createShortURL(t *testing.T) {
	usecase := usecase.New(repo.New())
	handler := chi.NewRouter()
	NewShortenRoutes(handler, usecase, logger.New("error"))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	cases := []struct {
		name                string
		body                string
		contentType         string
		expectedCode        int
		expectedContentType string
	}{
		{
			name:         "method_post_bad_content_type",
			body:         `{"url": "http://example.com"}`,
			contentType:  "text/plain; charset=utf-8",
			expectedCode: http.StatusUnsupportedMediaType,
		},
		{
			name:         "method_post_without_body",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "method_post_invalid_body_data",
			body:         `{}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:                "method_post_success",
			body:                `{"url": "http://example.com"}`,
			expectedCode:        http.StatusCreated,
			expectedContentType: "application/json",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			contentType := "application/json"

			if tt.contentType != "" {
				contentType = tt.contentType
			}

			resp, err := resty.New().
				R().
				SetHeader("Content-Type", contentType).
				SetBody(tt.body).
				Post(srv.URL + "/api/shorten")

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "Response code didn't match expected")

			if tt.expectedContentType != "" {
				assert.Equal(t, tt.expectedContentType, resp.Header().Get("Content-Type"))
			}
		})
	}
}
