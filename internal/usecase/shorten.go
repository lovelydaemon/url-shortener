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

func (uc *ShortenUseCase) Get(ctx context.Context, token string) (entity.StorageItem, error) {
	si, err := uc.repo.Get(ctx, token)
	if err != nil {
		return si, fmt.Errorf("ShortURLUseCase - Get - uc.repo.Get: %w", err)
	}
	return si, nil
}

func (uc *ShortenUseCase) Store(ctx context.Context, originalURL, token string) error {
	if err := uc.repo.Store(ctx, originalURL, token); err != nil {
		return fmt.Errorf("ShortURLUseCase - Store - uc.repo.Store: %w", err)
	}
	return nil
}
