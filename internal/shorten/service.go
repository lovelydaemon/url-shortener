package shorten

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lovelydaemon/url-shortener/internal/entity"
	"github.com/lovelydaemon/url-shortener/internal/repository"
)

type Service interface {
	Get(ctx context.Context, shortURL string) (entity.StorageURL, error)
	Store(ctx context.Context, originalURL entity.URL) (shortURL string, err error)
	StoreBatch(ctx context.Context, batch []BatchRequest) ([]BatchResponse, error)
}

type service struct {
	repo repository.Shorten
}

// NewService creates a new shorten service
func NewService(repo repository.Shorten) Service {
	return service{
		repo: repo,
	}
}

func (s service) Get(ctx context.Context, shortURL string) (entity.StorageURL, error) {
	storageURL, err := s.repo.Get(ctx, shortURL)
	if err != nil {
		return storageURL, fmt.Errorf("Shorten - service - Get - s.repo.Get: %w", err)
	}
	return storageURL, nil
}

func (s service) Store(ctx context.Context, originalURL entity.URL) (string, error) {
	shortURL, err := s.repo.Store(ctx, originalURL)
	if err != nil {
		return shortURL, fmt.Errorf("Shorten - service - Store - s.repo.Store: %w", err)
	}
	return shortURL, nil
}

type BatchRequest struct {
	ID          uuid.UUID  `json:"correlation_id"`
	OriginalURL entity.URL `json:"original_url"`
}

type BatchResponse struct {
	ID       uuid.UUID  `json:"correlation_id"`
	ShortURL entity.URL `json:"short_url"`
}

func (s service) StoreBatch(ctx context.Context, batch []BatchRequest) ([]BatchResponse, error) {
	storageUrls := make([]entity.StorageURL, 0, len(batch))
	for _, item := range batch {
		storageURL := entity.StorageURL{ID: item.ID, OriginalURL: item.OriginalURL}
		storageUrls = append(storageUrls, storageURL)
	}

	response, err := s.repo.StoreBatch(ctx, storageUrls)
	if err != nil {
		return nil, fmt.Errorf("Shorten - service - StoreBatch - s.repo.StoreBatch: %w", err)
	}

	batchResponse := make([]BatchResponse, 0, len(response))
	for _, item := range response {
		batchResponse = append(batchResponse, BatchResponse{ID: item.ID, ShortURL: entity.URL(item.ShortURL)})
	}

	return batchResponse, nil
}
