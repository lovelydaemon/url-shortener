package inmemory

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lovelydaemon/url-shortener/internal/entity"
	e "github.com/lovelydaemon/url-shortener/internal/pkg/errors"
	"github.com/lovelydaemon/url-shortener/internal/pkg/random"
)

func (s *storage) Get(ctx context.Context, shortURL string) (entity.StorageURL, error) {
	s.mu.RLock()
	storageURL, ok := s.storage[shortURL]
	s.mu.RUnlock()
	if !ok {
		return storageURL, e.ErrRecNotFound
	}
	return storageURL, nil
}

func (s *storage) Store(ctx context.Context, originalURL entity.URL) (string, error) {
	userID := ctx.Value("userID").(uuid.UUID)

	s.mu.RLock()
	for _, item := range s.storage {
		if item.OriginalURL == originalURL {
			s.mu.RUnlock()
			return item.ShortURL, e.ErrConflict
		}
	}
	s.mu.RUnlock()

	shortURL := random.NewRandomString()

	storageURL := entity.StorageURL{
		ID:          uuid.New(),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		UserID:      userID,
		CreatedAt:   time.Now(),
		DeletedFlag: false,
	}

	s.mu.Lock()
	s.storage[shortURL] = storageURL
	s.mu.Unlock()

	return shortURL, nil
}

func (s *storage) StoreBatch(ctx context.Context, batch []entity.StorageURL) ([]entity.StorageURL, error) {
	userID := ctx.Value("userID").(uuid.UUID)

	response := make([]entity.StorageURL, 0, len(batch))

	for _, item := range batch {
		shortURL := random.NewRandomString()

		response = append(response, entity.StorageURL{ID: item.ID, ShortURL: shortURL})

		storageURL := entity.StorageURL{
			ID:          uuid.New(),
			ShortURL:    shortURL,
			OriginalURL: item.OriginalURL,
			UserID:      userID,
			CreatedAt:   time.Now(),
			DeletedFlag: false,
		}

		s.mu.Lock()
		s.storage[shortURL] = storageURL
		s.mu.Unlock()
	}

	return response, nil
}
