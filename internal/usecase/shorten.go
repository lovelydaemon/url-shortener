package usecase

import (
	"context"
	"fmt"

	"github.com/lovelydaemon/url-shortener/internal/entity"
)

type ShortenUseCase struct {
	repo ShortenRepo
}

func NewShorten(r ShortenRepo) *ShortenUseCase {
	return &ShortenUseCase{
		repo: r,
	}
}

func (uc *ShortenUseCase) Get(ctx context.Context, shortURL string) (entity.Storage, error) {
	si, err := uc.repo.Get(ctx, shortURL)
	if err != nil {
		return si, fmt.Errorf("ShortenUseCase - Get - uc.repo.Get: %w", err)
	}
	return si, nil
}

func (uc *ShortenUseCase) Store(ctx context.Context, originalURL string) (string, error) {
	shortURL, err := uc.repo.Store(ctx, originalURL)
	if err != nil {
		return shortURL, fmt.Errorf("ShortenUseCase - Store - uc.repo.Store: %w", err)
	}
	return shortURL, nil
}

func (uc *ShortenUseCase) StoreBatch(ctx context.Context, batch []entity.BatchItemIn) ([]entity.BatchItemOut, error) {
	resp, err := uc.repo.StoreBatch(ctx, batch)
	if err != nil {
		return nil, fmt.Errorf("ShortenUseCase - StoreBatch - uc.repo.StoreBatch: %w", err)
	}
	return resp, nil
}
