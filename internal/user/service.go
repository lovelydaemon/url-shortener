package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lovelydaemon/url-shortener/internal/entity"
	"github.com/lovelydaemon/url-shortener/internal/queue"
	"github.com/lovelydaemon/url-shortener/internal/repository"
)

type Service interface {
	// GetURLs returns user urls
	GetURLs(ctx context.Context) (getURLsResponse, error)
	// DeleteURLs sends request to the deletion queue
	DeleteURLs(ctx context.Context, shortURLs []string)
}

type service struct {
	repo  repository.User
	queue *queue.Queue
}

// NewService creates a new user service
func NewService(repo repository.User, queue *queue.Queue) Service {
	return service{
		repo:  repo,
		queue: queue,
	}
}

type getURLsResponse []entity.UserURL

// GetURLs returns user urls
func (s service) GetURLs(ctx context.Context) (getURLsResponse, error) {
	response, err := s.repo.GetURLs(ctx)
	if err != nil {
		return nil, fmt.Errorf("User - service - GetURLs - s.repo.GetURLs: %w", err)
	}
	return response, nil
}

// DeleteURLs sends request to the deletion queue
func (s service) DeleteURLs(ctx context.Context, shortURLs []string) {
	userID := ctx.Value("userID").(uuid.UUID)

	u := entity.UserWithURLs{
		UserID:    userID,
		ShortURLs: shortURLs,
	}
	s.queue.Push(u)
}
