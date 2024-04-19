package usecase

import (
	"context"

	"github.com/lovelydaemon/url-shortener/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks.go -package=usecase

type (
	Shorten interface {
		Get(ctx context.Context, shortURL string) (entity.Storage, error)
		Store(ctx context.Context, originalUrl string) (string, error)
		StoreBatch(ctx context.Context, batch []entity.BatchItemIn) ([]entity.BatchItemOut, error)
	}

	ShortenRepo interface {
		Get(ctx context.Context, shortURL string) (entity.Storage, error)
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

type (
	User interface {
		GetURLs(ctx context.Context) ([]entity.UserURL, error)
		DeleteURLs(ctx context.Context, urls []string)
	}

	UserRepo interface {
		GetURLs(ctx context.Context) ([]entity.UserURL, error)
		DeleteURLs(ctx context.Context, urls ...entity.StorageWithUser) error
	}
)
