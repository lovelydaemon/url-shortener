package v1

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
	"github.com/lovelydaemon/url-shortener/internal/usecase/repo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ShortURLRoutes_getOriginalURL(t *testing.T) {
	usecase := usecase.New(repo.New())

	srv := httptest.NewServer(NewShortURLRoutes(usecase, "", logger.New("error")))
	defer srv.Close()

	originalURL := "http://example.com"
	client := resty.New().R()
	client.Method = http.MethodPost
	client.URL = srv.URL
	client.SetHeader("Content-type", "text/plain; charset=utf-8")
	client.SetBody(originalURL)
	resp, err := client.Send()
	require.NoError(t, err, "error making HTTP request")

	respURL := string(resp.Body())

	cases := []struct {
		name         string
		url          string
		expectedCode int
	}{
		{
			name:         "method_get_success_redirect",
			url:          respURL,
			expectedCode: http.StatusTemporaryRedirect,
		},
		{
			name:         "method_get_not_found",
			url:          fmt.Sprintf("%s/asdf", srv.URL),
			expectedCode: http.StatusNotFound,
		},
	}

	errRedirectBlocked := errors.New("HTTP redirect blocked")
	redirPolicy := resty.RedirectPolicyFunc(func(_ *http.Request, _ []*http.Request) error {
		return errRedirectBlocked
	})

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().
				SetRedirectPolicy(redirPolicy).
				R()
			resp, err := req.Get(tt.url)

			if !errors.Is(err, errRedirectBlocked) {
				assert.NoError(t, err, "error making HTTP request")
			}

			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "Response code didn't match expected")

			if resp.StatusCode() == http.StatusTemporaryRedirect {
				assert.Equal(t, originalURL, resp.Header().Get("Location"), "Location address didn't match expected")
			}
		})
	}
}

func Test_shortURLRoutes_createShortURL(t *testing.T) {
	usecase := usecase.New(repo.New())

	baseURL := "http://localhost:1234"

	srv := httptest.NewServer(NewShortURLRoutes(usecase, baseURL, logger.New("error")))
	defer srv.Close()

	cases := []struct {
		name                string
		bodyURL             string
		contentType         string
		expectedCode        int
		expectedResponseURL bool
	}{
		{
			name:         "method_post_bad_content_type",
			bodyURL:      "http://example.com",
			contentType:  "application/json",
			expectedCode: http.StatusUnsupportedMediaType,
		},
		{
			name:         "method_post_empty_body",
			bodyURL:      "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "method_post_bad_body_data",
			bodyURL:      "example.com",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:                "method_post_success",
			bodyURL:             "https://example.com",
			expectedCode:        http.StatusCreated,
			expectedResponseURL: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = http.MethodPost
			req.URL = srv.URL

			if tt.contentType != "" {
				req.SetHeader("Content-Type", tt.contentType)
			} else {
				req.SetHeader("Content-Type", "text/plain; charset=utf-8")
			}

			req.SetBody(tt.bodyURL)

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "Response code didn't match expected")

			if tt.expectedResponseURL {
				assert.True(t, strings.HasPrefix(string(resp.Body()), baseURL), "Response url prefix didn't match")
			}
		})
	}
}
