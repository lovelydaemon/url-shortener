package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lovelydaemon/url-shortener/internal/entity"
)

// может быть сделать хранилище обновление вот тут
type UserUseCase struct {
	repo      UserRepo
	storageCh chan entity.StorageWithUser
}

func NewUser(r UserRepo) *UserUseCase {
	u := &UserUseCase{
		repo:      r,
		storageCh: make(chan entity.StorageWithUser, 1024),
	}

	go u.flushStorage()

	return u
}

func (uc *UserUseCase) flushStorage() {
	ticker := time.NewTicker(time.Second * 10)

	var urls []entity.StorageWithUser

	for {
		select {
		case data := <-uc.storageCh:
			urls = append(urls, data)
		case <-ticker.C:
			if len(urls) == 0 {
				continue
			}

			if err := uc.repo.DeleteURLs(context.TODO(), urls...); err != nil {
				log.Printf("UserUseCase - flushStorage - uc.repo.DeleteURLs: %s", err.Error())
				continue
			}

			urls = nil
		}
	}
}

func (uc *UserUseCase) GetURLs(ctx context.Context) ([]entity.UserURL, error) {
	userURLs, err := uc.repo.GetURLs(ctx)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - GetURLs - uc.repo.GetUrls: %w", err)
	}
	return userURLs, nil
}

func (uc *UserUseCase) DeleteURLs(ctx context.Context, urls []string) {
	go func() {
		userID := ctx.Value("userID").(uuid.UUID)

		for _, url := range urls {
			data := entity.StorageWithUser{
				UserID:   userID,
				ShortURL: url,
			}
			select {
			case <-ctx.Done():
				return
			case uc.storageCh <- data:
			}
		}
	}()
}
