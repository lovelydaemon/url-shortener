package repo

import (
	"context"
	"fmt"

	"github.com/lovelydaemon/url-shortener/internal/entity"
	"github.com/lovelydaemon/url-shortener/internal/postgres"
)

type ShortenRepoPG struct {
	*postgres.Postgres
}

func NewShortenPG(pg *postgres.Postgres) *ShortenRepoPG {
	return &ShortenRepoPG{
		pg,
	}
}

func (r *ShortenRepoPG) Get(ctx context.Context, token string) (entity.StorageItem, error) {
	var si entity.StorageItem

	err := r.Pool.QueryRow(ctx,
		"SELECT id, short_url, original_url FROM urls WHERE short_url = $1", token).
		Scan(&si.ID, &si.Token, &si.OriginalURL)

	if err != nil {
		return si, fmt.Errorf("ShortURLRepo - Get - r.Pool.QueryRow.Scan: %w", err)
	}

	return si, nil
}

func (r *ShortenRepoPG) Store(ctx context.Context, originalURL, token string) error {
	_, err := r.Pool.Exec(ctx,
		"INSERT INTO urls (short_url, original_url) VALUES ($1, $2)", token, originalURL)
	if err != nil {
		return fmt.Errorf("ShortURLRepo - Store - r.Pool.Exec: %w", err)
	}
	return nil
}
