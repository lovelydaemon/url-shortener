package usecase

import (
	"context"

	"github.com/lovelydaemon/url-shortener/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks.go -package=usecase

type (
	Shorten interface {
		Get(ctx context.Context, token string) (entity.StorageItem, error)
		Store(ctx context.Context, originalUrl string) (string, error)
	}

	ShortenRepo interface {
		Get(ctx context.Context, token string) (entity.StorageItem, error)
		Store(ctx context.Context, originalUrl, token string) error
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
