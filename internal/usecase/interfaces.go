package usecase

import "github.com/lovelydaemon/url-shortener/internal/entity"

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
