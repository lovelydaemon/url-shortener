package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var errInternalServerErr = errors.New("internal server error")

func ping(t *testing.T) (*PingUseCase, *MockPingRepo) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockPingRepo(ctrl)
	ping := NewPing(repo)

	return ping, repo
}

func TestPing(t *testing.T) {
	ping, repo := ping(t)

	tests := []struct {
		name string
		mock func()
		err  error
	}{
		{
			name: "success",
			mock: func() {
				repo.EXPECT().Ping(context.Background()).Return(nil)
			},
			err: nil,
		},
		{
			name: "error",
			mock: func() {
				repo.EXPECT().Ping(context.Background()).Return(errInternalServerErr)
			},
			err: errInternalServerErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := ping.Ping(context.Background())
			assert.ErrorIs(t, err, tt.err)
		})
	}

}
