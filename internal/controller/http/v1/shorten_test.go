package v1

import (
	"fmt"
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

func shorten(t *testing.T) (*httptest.Server, *usecase.MockShortenRepo) {
	t.Helper()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := usecase.NewMockShortenRepo(ctrl)
	uc := usecase.NewShorten(repo)

	handler := chi.NewRouter()
	NewShortenRoutes(handler, logger.New("error"), uc)
	srv := httptest.NewServer(handler)

	return srv, repo
}

func Test_shortenRoutes_generateShortURL(t *testing.T) {
	srv, repo := shorten(t)
	defer srv.Close()

	tests := []struct {
		name                string
		mock                func()
		contentType         string
		body                string
		expectedBody        string
		expectedContentType string
		expectedCode        int
	}{
		{
			name:         "bad content type",
			mock:         func() {},
			contentType:  "text/plain",
			body:         ``,
			expectedCode: http.StatusUnsupportedMediaType,
		},
		{
			name:         "error decoding body",
			mock:         func() {},
			body:         `{"url: http://example.com,}`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "bad request body url",
			mock:         func() {},
			body:         `{"url": "example.com"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "error on store",
			mock: func() {
				repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return("", ErrInternalServerError)
			},
			body:         `{"url": "http://example.com"}`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "conflict url already exists",
			mock: func() {
				repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return("abc", ErrConflict)
			},
			body:                `{"url": "http://example.com"}`,
			expectedContentType: "application/json",
			expectedBody:        fmt.Sprintf("{\"result\":\"%s/abc\"}\n", srv.URL),
			expectedCode:        http.StatusConflict,
		},
		{
			name: "success",
			mock: func() {
				repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return("abc", nil)
			},
			body:                `{"url": "http://example.com"}`,
			expectedContentType: "application/json",
			expectedBody:        fmt.Sprintf(`{"result":"%s/abc"}`, srv.URL),
			expectedCode:        http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			contentType := "application/json"
			if tt.contentType != "" {
				contentType = tt.contentType
			}

			resp, err := resty.New().
				R().
				SetBody(tt.body).
				SetHeader("Content-Type", contentType).
				Post(srv.URL + "/api/shorten")
			require.NoError(t, err, "Error making HTTP request")

			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "Response code didn't match expected")

			if tt.expectedContentType != "" {
				assert.Equal(t, "application/json", resp.Header().Get("Content-Type"),
					"Response Content-Type didn't match expected",
				)
			}

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, string(resp.Body()), "Response body didn't match expected")
			}
		})
	}
}

func Test_shortenRoutes_generateShortURLBatch(t *testing.T) {
	srv, repo := shorten(t)
	defer srv.Close()

	requestBody := `[
  {"correlation_id": "1", "original_url": "http://example.com"},
  {"correlation_id": "2", "original_url": "http://example2.com"}
  ]`

	successResponseBody := fmt.Sprintf(`[
  {"correlation_id": "1", "short_url": "%s/abc"},
  {"correlation_id": "2", "short_url": "%s/abcd"}
  ]`, srv.URL, srv.URL)

	tests := []struct {
		name                string
		mock                func()
		contentType         string
		body                string
		expectedBody        string
		expectedContentType string
		expectedCode        int
	}{
		{
			name:         "bad content type",
			mock:         func() {},
			contentType:  "text/plain",
			body:         ``,
			expectedCode: http.StatusUnsupportedMediaType,
		},
		{
			name:         "error decoding body",
			mock:         func() {},
			body:         "",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "empty request",
			mock:         func() {},
			body:         `[]`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "bad request body url",
			mock:         func() {},
			body:         `[{"correlation_id": "1", "original_url": "example.com"}]`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "error on store batch",
			mock: func() {
				repo.EXPECT().StoreBatch(gomock.Any(), gomock.Any()).Return(nil, ErrInternalServerError)
			},
			body:         requestBody,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "success",
			mock: func() {
				repo.EXPECT().StoreBatch(gomock.Any(), gomock.Any()).
					Return([]entity.BatchItemOut{
						{ID: "1", ShortURL: "abc"},
						{ID: "2", ShortURL: "abcd"},
					}, nil)
			},
			body:                requestBody,
			expectedContentType: "application/json",
			expectedBody:        successResponseBody,
			expectedCode:        http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			contentType := "application/json"
			if tt.contentType != "" {
				contentType = tt.contentType
			}

			resp, err := resty.New().
				R().
				SetBody(tt.body).
				SetHeader("Content-Type", contentType).
				Post(srv.URL + "/api/shorten/batch")
			require.NoError(t, err, "Error making HTTP request")

			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "Response code didn't match expected")

			if tt.expectedContentType != "" {
				assert.Equal(t, "application/json", resp.Header().Get("Content-Type"),
					"Response Content-Type didn't match expected",
				)
			}

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, string(resp.Body()), "Response body didn't match expected")
			}
		})
	}
}
