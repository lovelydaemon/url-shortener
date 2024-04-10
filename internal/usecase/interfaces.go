package usecase

import (
	"context"

	"github.com/lovelydaemon/url-shortener/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	Shorten interface {
		Get(ctx context.Context, token string) (entity.StorageItem, error)
		Store(ctx context.Context, originalUrl string) (string, error)
		StoreBatch(ctx context.Context, batch []entity.BatchItemIn) ([]entity.BatchItemOut, error)
	}

	ShortenRepo interface {
		Get(ctx context.Context, token string) (entity.StorageItem, error)
		Store(ctx context.Context, originalUrl string) (string, error)
		StoreBatch(ctx context.Context, batch []entity.BatchItemIn) ([]entity.BatchItemOut, error)
	}
)

type (
	Ping interface {
		Ping(ctx context.Context) error
	}

	PingRepo interface {
		Ping(ctx context.Context) error
	}
)
