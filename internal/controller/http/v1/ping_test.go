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
	"go.uber.org/mock/gomock"
)

func ping(t *testing.T) (*usecase.PingUseCase, *usecase.MockPingRepo) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := usecase.NewMockPingRepo(ctrl)
	ping := usecase.NewPing(repo)
	return ping, repo
}

func Test_Ping_ping(t *testing.T) {
	ping, repo := ping(t)

	handler := chi.NewRouter()
	NewPingRoutes(handler, logger.New("error"), ping)
	srv := httptest.NewServer(handler)
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
			name: "server_error",
			mock: func() {
				repo.EXPECT().Ping(gomock.Any()).Return(errInternalServerError)
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			resp, err := resty.New().R().Get(srv.URL + "/ping")
			assert.NoError(t, err, "error making HTTP request")

			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "Response code didn't match expected")
		})
	}
}
