package repo

import (
	"context"
	"fmt"

	"github.com/lovelydaemon/url-shortener/internal/entity"
	"github.com/lovelydaemon/url-shortener/internal/random"
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
		return si, fmt.Errorf("ShortenRepo - Get - r.storage.Get: %w", err)
	}
	return si, nil
}

func (r *ShortenRepoST) Store(ctx context.Context, originalURL string) (string, error) {
	token := random.NewRandomString()

	storageItem := entity.StorageItem{
		ID:          r.storage.Len() + 1,
		Token:       token,
		OriginalURL: originalURL,
	}

	if err := r.storage.Write(storageItem); err != nil {
		return token, fmt.Errorf("ShortenRepo - Store - r.storage.Write: %w", err)
	}
	return token, nil
}

func (r *ShortenRepoST) StoreBatch(ctx context.Context, batch []entity.BatchItemIn) ([]entity.BatchItemOut, error) {
	return nil, nil
}
