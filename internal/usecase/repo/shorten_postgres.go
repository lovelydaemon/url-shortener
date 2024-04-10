package repo

import (
	"context"
	"fmt"

	"github.com/lovelydaemon/url-shortener/internal/entity"
	"github.com/lovelydaemon/url-shortener/internal/postgres"
	"github.com/lovelydaemon/url-shortener/internal/random"
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
		return si, fmt.Errorf("ShortenRepo - Get - r.Pool.QueryRow.Scan: %w", err)
	}

	return si, nil
}

func (r *ShortenRepoPG) Store(ctx context.Context, originalURL string) (string, error) {
	token := random.NewRandomString()

	_, err := r.Pool.Exec(ctx,
		"INSERT INTO urls (short_url, original_url) VALUES ($1, $2)", token, originalURL)
	if err != nil {
		return token, fmt.Errorf("ShortenRepo - Store - r.Pool.Exec: %w", err)
	}
	return token, nil
}

func (r *ShortenRepoPG) StoreBatch(ctx context.Context, batch []entity.BatchItemIn) ([]entity.BatchItemOut, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("ShortenRepo - StoreBatch - r.Pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)

	correlationIds := make([]string, 0, len(batch))
	for _, v := range batch {
		token := random.NewRandomString()
		correlationIds = append(correlationIds, v.ID)

		_, err = tx.Exec(ctx, `
      INSERT INTO urls (short_url, original_url, correlation_id)
      VALUES ($1, $2, $3)
    `, token, v.OriginalURL, v.ID)

		if err != nil {
			return nil, fmt.Errorf("ShortenRepo - StoreBatch - tx.Exec: %w", err)
		}
	}
	tx.Commit(ctx)

	rows, err := r.Pool.Query(ctx, `
    SELECT correlation_id, short_url FROM urls
    WHERE correlation_id = ANY($1)
  `, correlationIds)
	if err != nil {
		return nil, fmt.Errorf("ShortenRepo - StoreBatch - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var batchOut []entity.BatchItemOut
	for rows.Next() {
		var item entity.BatchItemOut
		if err := rows.Scan(&item.ID, &item.ShortURL); err != nil {
			return nil, fmt.Errorf("ShortenRepo - StoreBatch - rows.Scan: %w", err)
		}

		batchOut = append(batchOut, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("ShortenRepo - StoreBatch - rows.Err: %w", err)
	}

	return batchOut, nil
}
