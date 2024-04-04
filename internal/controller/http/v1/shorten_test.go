package v1

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/storage"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
	"github.com/lovelydaemon/url-shortener/internal/usecase/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ShortenRoutes_createShortURL(t *testing.T) {
	st, err := storage.NewStorage("")
	require.NoError(t, err, "Couldn't create storage")

	usecase := usecase.New(repo.New(st))
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
		expectedBody        string
	}{
		{
			name:         "method_post_bad_content_type",
			body:         `{"url": "http://example.com"}`,
			contentType:  "text/plain; charset=utf-8",
			expectedCode: http.StatusUnsupportedMediaType,
			expectedBody: "",
		},
		{
			name:         "method_post_without_body",
			expectedCode: http.StatusInternalServerError,
			expectedBody: "",
		},
		{
			name:         "method_post_invalid_body_data",
			body:         `{}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:                "method_post_success",
			body:                `{"url": "http://example.com"}`,
			expectedCode:        http.StatusCreated,
			expectedContentType: "application/json",
			expectedBody:        fmt.Sprintf("%s/.........", srv.URL),
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
				assert.Equal(t, tt.expectedContentType, resp.Header().Get("Content-Type"), "Response Content-Type didn't match expected")
			}

			if tt.expectedBody != "" {
				assert.Regexp(t, tt.expectedBody, string(resp.Body()), "Response url didn't match expected")
			}
		})
	}
}
