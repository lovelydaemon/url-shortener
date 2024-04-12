package v1

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/lovelydaemon/url-shortener/internal/entity"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func shortURL(t *testing.T, shortAddr string) (*httptest.Server, *usecase.MockShortenRepo) {
	t.Helper()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := usecase.NewMockShortenRepo(ctrl)
	uc := usecase.NewShorten(repo)

	handler := chi.NewRouter()
	NewShortURLRoutes(handler, logger.New("error"), uc, shortAddr)
	srv := httptest.NewServer(handler)

	return srv, repo
}

func Test_shortURLRoutes_getOriginalURL(t *testing.T) {
	shortAddr := ""
	originalURL := "http://example.com"
	srv, repo := shortURL(t, shortAddr)
	defer srv.Close()

	tests := []struct {
		name             string
		token            string
		mock             func()
		expectedCode     int
		expectedLocation string
	}{
		{
			name:  "success",
			token: "abc",
			mock: func() {
				repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(entity.StorageItem{OriginalURL: originalURL}, nil)
			},
			expectedCode:     http.StatusTemporaryRedirect,
			expectedLocation: originalURL,
		},
		{
			name:  "not found",
			token: "abcd",
			mock: func() {
				repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(entity.StorageItem{}, ErrInternalServerError)
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			resp, err := resty.
				New().
				SetRedirectPolicy(resty.NoRedirectPolicy()).
				R().
				Get(srv.URL + "/" + tt.token)

			if !errors.Is(err, resty.ErrAutoRedirectDisabled) {
				require.NoError(t, err, "Error making HTTP request")
			}

			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "Response code didn't match expected")

			if resp.StatusCode() == http.StatusTemporaryRedirect {
				assert.Equal(t, tt.expectedLocation, resp.Header().Get("Location"), "Location URL didn't match expected")
			}
		})
	}
}

func Test_shortURLRoutes_generateShortURL(t *testing.T) {
	srv, repo := shortURL(t, "")
	defer srv.Close()

	tests := []struct {
		name                string
		mock                func()
		body                string
		expectedContentType string
		expectedBody        string
		expectedCode        int
	}{
		{
			name: "success",
			mock: func() {
				repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return("abc", nil)
			},
			body:                "http://example.com",
			expectedContentType: "text/plain",
			expectedBody:        srv.URL + "/abc",
			expectedCode:        http.StatusCreated,
		},
		{
			name:         "bad request body url",
			mock:         func() {},
			body:         "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "error on store",
			mock: func() {
				repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return("", ErrInternalServerError)
			},
			body:         "http://example.com",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "url already exists",
			mock: func() {
				repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return("abc", ErrConflict)
			},
			body:                "http://example.com",
			expectedContentType: "text/plain",
			expectedBody:        srv.URL + "/abc",
			expectedCode:        http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			resp, err := resty.
				New().
				R().
				SetBody(tt.body).
				Post(srv.URL + "/")
			require.NoError(t, err, "Error making HTTP request")

			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "Response code didn't match expected")

			if tt.expectedContentType != "" {
				assert.Equal(
					t, tt.expectedContentType, resp.Header().Get("Content-Type"),
					"Content-Type didn't match expected",
				)
			}

			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, string(resp.Body()), "Response body didn't match expected")
			}
		})
	}
}

func Test_shortURLRoutes_generateShortURL_with_shortAddr(t *testing.T) {
	shortAddr := "http://example.com:1234"
	srv, repo := shortURL(t, shortAddr)
	defer srv.Close()

	tests := []struct {
		name         string
		mock         func()
		body         string
		expectedBody string
		expectedCode int
	}{
		{
			name: "success",
			mock: func() {
				repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return("abc", nil)
			},
			body:         "http://example.com",
			expectedBody: shortAddr + "/abc",
			expectedCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			resp, err := resty.
				New().
				R().
				SetBody(tt.body).
				Post(srv.URL + "/")
			require.NoError(t, err, "Error making HTTP request")

			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "Response code didn't match expected")
			assert.Equal(t, tt.expectedBody, string(resp.Body()), "Response body didn't match expected")
		})
	}
}
