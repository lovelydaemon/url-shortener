package usecase_test

import (
	"context"
	"testing"

	"github.com/lovelydaemon/url-shortener/internal/entity"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func shorten(t *testing.T) (*usecase.ShortenUseCase, *MockShortenRepo) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockShortenRepo(ctrl)
	shorten := usecase.NewShorten(repo)

	return shorten, repo
}

func TestGet(t *testing.T) {
	shorten, repo := shorten(t)

	tests := []struct {
		name string
		mock func()
		res  any
		err  error
	}{
		{
			name: "success",
			mock: func() {
				repo.EXPECT().Get(context.Background(), "").Return(entity.StorageItem{}, nil)
			},
			res: entity.StorageItem{},
			err: nil,
		},
		{
			name: "empty result with error",
			mock: func() {
				repo.EXPECT().Get(context.Background(), "").Return(entity.StorageItem{}, errInternalServerErr)
			},
			res: entity.StorageItem{},
			err: errInternalServerErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			res, err := shorten.Get(context.Background(), "")

			assert.Equal(t, tt.res, res)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func TestStore(t *testing.T) {
	shorten, repo := shorten(t)

	tests := []struct {
		name string
		mock func()
		res  any
		err  error
	}{
		{
			name: "success",
			mock: func() {
				repo.EXPECT().Store(context.Background(), "").Return("token", nil)
			},
			res: "token",
			err: nil,
		},
		{
			name: "empty result with error",
			mock: func() {
				repo.EXPECT().Store(context.Background(), "").Return("", errInternalServerErr)
			},
			res: "",
			err: errInternalServerErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			res, err := shorten.Store(context.Background(), "")

			assert.Equal(t, tt.res, res)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func TestStoreBatch(t *testing.T) {
	shorten, repo := shorten(t)

	tests := []struct {
		name string
		mock func()
		res  any
		err  error
	}{
		{
			name: "success",
			mock: func() {
				repo.EXPECT().StoreBatch(context.Background(), []entity.BatchItemIn{}).Return([]entity.BatchItemOut{}, nil)
			},
			res: []entity.BatchItemOut{},
			err: nil,
		},
		{
			name: "empty result with error",
			mock: func() {
				repo.EXPECT().StoreBatch(context.Background(), []entity.BatchItemIn{}).Return(nil, errInternalServerErr)
			},
			res: []entity.BatchItemOut(nil),
			err: errInternalServerErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			res, err := shorten.StoreBatch(context.Background(), []entity.BatchItemIn{})

			assert.Equal(t, tt.res, res)
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
