package repository

import (
	"context"

	"github.com/lovelydaemon/url-shortener/internal/entity"
)

type (
	Shorten interface {
		Get(ctx context.Context, shortURL string) (entity.StorageURL, error)
		Store(ctx context.Context, originalURL entity.URL) (shortURL string, err error)
		StoreBatch(ctx context.Context, batch []entity.StorageURL) ([]entity.StorageURL, error)
	}

	Ping interface {
		Ping(ctx context.Context) error
	}

	User interface {
		GetURLs(ctx context.Context) ([]entity.UserURL, error)
		DeleteURLs(ctx context.Context, shortURLs ...entity.UserWithURLs) error
	}
)
