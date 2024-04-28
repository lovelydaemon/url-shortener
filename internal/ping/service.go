package ping

import (
	"context"
	"fmt"

	"github.com/lovelydaemon/url-shortener/internal/repository"
)

type Service interface {
	Ping(ctx context.Context) error
}

type service struct {
	repo repository.Ping
}

// NewService creates a new ping service
func NewService(repo repository.Ping) Service {
	return service{
		repo: repo,
	}
}

func (s service) Ping(ctx context.Context) error {
	if err := s.repo.Ping(ctx); err != nil {
		return fmt.Errorf("Ping - service - Ping - s.repo.Ping: %w", err)
	}
	return nil
}
