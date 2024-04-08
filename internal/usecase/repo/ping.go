package repo

import (
	"context"
	"fmt"

	"github.com/lovelydaemon/url-shortener/internal/postgres"
)

type PingRepo struct {
	*postgres.Postgres
}

func NewPing(pg *postgres.Postgres) *PingRepo {
	return &PingRepo{pg}
}

func (r *PingRepo) Ping(ctx context.Context) error {
	if err := r.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("PingRepo - Ping - r.Pool.Ping: %w", err)
	}
	return nil
}
