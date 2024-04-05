package repo

import (
	"github.com/lovelydaemon/url-shortener/internal/entity"
	"github.com/lovelydaemon/url-shortener/internal/storage"
)

type ShortURLRepo struct {
	storage *storage.Storage
}

func NewShortURLRepo(storage *storage.Storage) *ShortURLRepo {
	return &ShortURLRepo{
		storage: storage,
	}
}

func (r *ShortURLRepo) Get(token string) (entity.StorageItem, bool) {
	for _, v := range r.storage.Store {
		if v.Token == token {
			return v, true
		}
	}

	return entity.StorageItem{}, false
}

func (r *ShortURLRepo) Store(originalURL, token string) error {
	uuid := len(r.storage.Store) + 1

	storageItem := entity.StorageItem{
		UUID:        uuid,
		Token:       token,
		OriginalURL: originalURL,
	}

	return r.storage.Write(storageItem)
}
