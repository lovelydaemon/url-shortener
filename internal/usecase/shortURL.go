package usecase

import (
	"fmt"

	"github.com/lovelydaemon/url-shortener/internal/entity"
)

type ShortURLUseCase struct {
	repo ShortURLRepo
}

func New(r ShortURLRepo) *ShortURLUseCase {
	return &ShortURLUseCase{
		repo: r,
	}
}

func (uc *ShortURLUseCase) Get(token string) (entity.StorageItem, bool) {
	return uc.repo.Get(token)
}

func (uc *ShortURLUseCase) Store(originalURL, token string) error {
	if err := uc.repo.Store(originalURL, token); err != nil {
		return fmt.Errorf("ShortURLUseCase - Store - uc.repo.Store: %w", err)
	}
	return nil
}
