package repo

import (
	"context"

	"github.com/lovelydaemon/url-shortener/internal/postgres"
)

type PingRepo struct {
	*postgres.Postgres
}

func NewPingRepo(pg *postgres.Postgres) *PingRepo {
	return &PingRepo{pg}
}

func (r *PingRepo) Ping(ctx context.Context) error {
	return r.Pool.Ping(ctx)
}
