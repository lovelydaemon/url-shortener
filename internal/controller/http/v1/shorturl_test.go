package v1

import (
	"errors"
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

func Test_ShortURLRoutes_getOriginalURL(t *testing.T) {
	st, err := storage.New("")
	require.NoError(t, err, "Couldn't create storage")

	usecase := usecase.NewShortURLUseCase(repo.NewShortURLRepo(st))
	handler := chi.NewRouter()
	NewShortURLRoutes(handler, logger.New("error"), usecase, "")
	srv := httptest.NewServer(handler)
	defer srv.Close()

	originalURL := "http://example.com"

	resp, err := resty.New().
		R().
		SetBody(originalURL).
		Post(srv.URL)

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
			resp, err := resty.New().
				SetRedirectPolicy(redirPolicy).
				R().
				Get(tt.url)

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
	st, err := storage.New("")
	require.NoError(t, err, "Couldn't create storage")

	usecase := usecase.NewShortURLUseCase(repo.NewShortURLRepo(st))
	handler := chi.NewRouter()
	NewShortURLRoutes(handler, logger.New("error"), usecase, "")
	srv := httptest.NewServer(handler)
	defer srv.Close()

	cases := []struct {
		name         string
		bodyURL      string
		contentType  string
		expectedCode int
		expectedBody string
	}{
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
			name:         "method_post_success",
			bodyURL:      "https://example.com",
			expectedCode: http.StatusCreated,
			expectedBody: fmt.Sprintf("%s/.........", srv.URL),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := resty.New().
				R().
				SetBody(tt.bodyURL).
				Post(srv.URL)

			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "Response code didn't match expected")

			if tt.expectedBody != "" {
				assert.Regexp(t, tt.expectedBody, string(resp.Body()), "Response url didn't match expected")
			}
		})
	}
}
