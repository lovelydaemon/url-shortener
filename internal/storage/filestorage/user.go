package filestorage

import (
	"context"

	"github.com/google/uuid"
	"github.com/lovelydaemon/url-shortener/internal/entity"
)

func (s *storage) GetURLs(ctx context.Context) ([]entity.UserURL, error) {
	userID := ctx.Value("userID").(uuid.UUID)

	var response []entity.UserURL

	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.ReadFromFile(); err != nil {
		return nil, err
	}

	for _, v := range s.storage {
		if v.UserID == userID {
			userURL := entity.UserURL{
				ShortURL:    entity.URL(v.ShortURL),
				OriginalURL: v.OriginalURL,
			}

			response = append(response, userURL)
		}
	}

	return response, nil
}

func (s *storage) DeleteURLs(ctx context.Context, data ...entity.UserWithURLs) error {
	s.mu.RLock()
	if err := s.ReadFromFile(); err != nil {
		return err
	}
	s.mu.RUnlock()

	for _, item := range data {
		for _, shortURL := range item.ShortURLs {
			s.mu.RLock()
			storageURL, ok := s.storage[shortURL]
			s.mu.RUnlock()
			if !ok || storageURL.UserID != item.UserID || storageURL.DeletedFlag {
				continue
			}

			storageURL.DeletedFlag = true

			s.mu.Lock()
			s.storage[shortURL] = storageURL
			s.mu.Unlock()
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.WriteToFile(); err != nil {
		return err
	}

	return nil
}
