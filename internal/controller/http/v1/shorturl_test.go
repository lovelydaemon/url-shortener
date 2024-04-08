package v1

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/lovelydaemon/url-shortener/internal/entity"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"
)

func Test_ShortURLRoutes_getOriginalURL(t *testing.T) {
	shorten, repo := shorten(t)

	handler := chi.NewRouter()
	NewShortURLRoutes(handler, logger.New("error"), shorten, "")
	srv := httptest.NewServer(handler)
	defer srv.Close()

	originalURL := "http://example.com"
	token := "abcdefg"

	tests := []struct {
		name         string
		url          string
		mock         func()
		expectedCode int
	}{
		{
			name: "success_redirect",
			url:  fmt.Sprintf("%s/%s", srv.URL, token),
			mock: func() {
				repo.EXPECT().Get(gomock.Any(), token).Return(entity.StorageItem{OriginalURL: originalURL}, nil)
			},
			expectedCode: http.StatusTemporaryRedirect,
		},
		{
			name: "not_found",
			url:  fmt.Sprintf("%s/%s", srv.URL, "abc"),
			mock: func() {
				repo.EXPECT().Get(gomock.Any(), "abc").Return(entity.StorageItem{}, errNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	redirPolicy := resty.RedirectPolicyFunc(func(_ *http.Request, _ []*http.Request) error {
		return errRedirectBlocked
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

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
	shorten, repo := shorten(t)

	handler := chi.NewRouter()
	NewShortURLRoutes(handler, logger.New("error"), shorten, "")
	srv := httptest.NewServer(handler)
	defer srv.Close()

	tests := []struct {
		name         string
		bodyURL      string
		contentType  string
		mock         func()
		expectedCode int
		expectedBody string
	}{
		{
			name:    "empty_body",
			bodyURL: "",
			mock: func() {
				repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "bad_body_data",
			bodyURL: "example.com",
			mock: func() {
				repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "success",
			bodyURL: "https://example.com",
			mock: func() {
				repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: fmt.Sprintf("%s/.........", srv.URL),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

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
