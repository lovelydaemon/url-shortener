package usecase

import (
	"context"

	"github.com/lovelydaemon/url-shortener/internal/entity"
)

type (
	ShortURL interface {
		Get(token string) (entity.StorageItem, bool)
		Store(originalUrl, token string) error
	}

	ShortURLRepo interface {
		Get(token string) (entity.StorageItem, bool)
		Store(originalUrl, token string) error
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
