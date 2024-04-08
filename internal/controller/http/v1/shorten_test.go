package v1

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
	"go.uber.org/mock/gomock"

	"github.com/stretchr/testify/assert"
)

func shorten(t *testing.T) (*usecase.ShortenUseCase, *usecase.MockShortenRepo) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := usecase.NewMockShortenRepo(ctrl)
	shorten := usecase.NewShorten(repo)
	return shorten, repo
}

func Test_ShortenRoutes_createShortURL(t *testing.T) {
	shorten, repo := shorten(t)

	handler := chi.NewRouter()
	NewShortenRoutes(handler, logger.New("error"), shorten)

	srv := httptest.NewServer(handler)
	defer srv.Close()

	tests := []struct {
		name                string
		body                string
		contentType         string
		mock                func()
		expectedCode        int
		expectedContentType string
		expectedBody        string
	}{
		{
			name:        "method_post_bad_content_type",
			body:        `{"url": "http://example.com"}`,
			contentType: "text/plain; charset=utf-8",
			mock: func() {
				repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusUnsupportedMediaType,
		},
		{
			name: "method_post_without_body",
			mock: func() {
				repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "method_post_invalid_body_data",
			body: `{}`,
			mock: func() {
				repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "method_post_success",
			body: `{"url": "http://example.com"}`,
			mock: func() {
				repo.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedCode:        http.StatusCreated,
			expectedContentType: "application/json",
			expectedBody:        fmt.Sprintf("%s/.........", srv.URL),
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
