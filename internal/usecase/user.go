package usecase

import (
	"context"
	"fmt"

	"github.com/lovelydaemon/url-shortener/internal/entity"
)

type UserUseCase struct {
	repo UserRepo
}

func NewUser(r UserRepo) *UserUseCase {
	return &UserUseCase{r}
}

func (uc *UserUseCase) GetUrls(ctx context.Context) ([]entity.UserURL, error) {
	userURLs, err := uc.repo.GetUrls(ctx)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - GetUrls - uc.repo.GetUrls: %w", err)
	}
	return userURLs, nil
}
