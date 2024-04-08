package repo

import (
	"context"
	"fmt"

	"github.com/lovelydaemon/url-shortener/internal/entity"
	"github.com/lovelydaemon/url-shortener/internal/storage"
)

type ShortenRepoST struct {
	storage *storage.Storage
}

func NewShortenST(storage *storage.Storage) *ShortenRepoST {
	return &ShortenRepoST{
		storage: storage,
	}
}

func (r *ShortenRepoST) Get(ctx context.Context, token string) (entity.StorageItem, error) {
	si, err := r.storage.Get(token)
	if err != nil {
		return si, fmt.Errorf("ShortenRepoST - Get - r.storage.Get: %w", err)
	}
	return si, nil
}

func (r *ShortenRepoST) Store(ctx context.Context, originalURL, token string) error {
	storageItem := entity.StorageItem{
		ID:          r.storage.Len() + 1,
		Token:       token,
		OriginalURL: originalURL,
	}

	if err := r.storage.Write(storageItem); err != nil {
		return fmt.Errorf("ShortenRepoST - Store - r.storage.Write: %w", err)
	}
	return nil
}
