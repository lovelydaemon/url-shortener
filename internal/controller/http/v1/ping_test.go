package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func ping(t *testing.T) (*httptest.Server, *usecase.MockPingRepo) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := usecase.NewMockPingRepo(ctrl)
	uc := usecase.NewPing(repo)
	handler := chi.NewRouter()
	NewPingRoutes(handler, logger.New("error"), uc)

	srv := httptest.NewServer(handler)

	return srv, repo
}

func Test_pingRoutes_ping(t *testing.T) {
	srv, repo := ping(t)
	defer srv.Close()

	tests := []struct {
		name         string
		mock         func()
		expectedCode int
	}{
		{
			name: "success",
			mock: func() {
				repo.EXPECT().Ping(gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "error ping db",
			mock: func() {
				repo.EXPECT().Ping(gomock.Any()).Return(ErrInternalServerError)
			},
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			resp, err := resty.
				New().
				R().
				Get(srv.URL + "/ping")
			require.NoError(t, err, "Error making HTTP request")

			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}
