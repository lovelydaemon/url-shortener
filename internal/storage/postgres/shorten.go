package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lovelydaemon/url-shortener/internal/entity"
	e "github.com/lovelydaemon/url-shortener/internal/pkg/errors"
	"github.com/lovelydaemon/url-shortener/internal/pkg/random"
)

func (s storage) Get(ctx context.Context, shortURL string) (entity.StorageURL, error) {
	var storageURL entity.StorageURL

	query := `
  SELECT original_url, is_deleted FROM urls
  WHERE short_url = $1
  `

	err := s.Pool.QueryRow(ctx, query, shortURL).
		Scan(&storageURL.OriginalURL, &storageURL.DeletedFlag)

	if err == nil {
		return storageURL, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return storageURL, fmt.Errorf("Storage - Get - s.Pool.QueryRow: %w", e.ErrRecNotFound)
	}

	return storageURL, fmt.Errorf("Storage - Get - s.Pool.QueryRow: %w", err)
}

func (s storage) Store(ctx context.Context, originalURL entity.URL) (string, error) {
	userID := ctx.Value("userID").(uuid.UUID)
	shortURL := ""

	query := `
  SELECT short_url FROM urls
  WHERE original_url = $1
  `

	err := s.Pool.QueryRow(ctx, query, originalURL).Scan(&shortURL)
	if err == nil {
		return shortURL, e.ErrConflict
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("Storage - Store - s.Pool.QueryRow: %w", err)
	}

	query = `
  INSERT INTO urls
  (short_url, original_url, user_id)
  VALUES
  ($1, $2, $3)
  `

	shortURL = random.NewRandomString()

	if _, err = s.Pool.Exec(ctx, query, shortURL, originalURL, userID); err != nil {
		return "", fmt.Errorf("Storage - Store - s.Pool.Exec: %w", err)
	}

	return shortURL, nil
}

func (s storage) StoreBatch(ctx context.Context, batch []entity.StorageURL) ([]entity.StorageURL, error) {
	userID := ctx.Value("userID").(uuid.UUID)

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("Storage - StoreBatch - s.Pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)

	txQuery := `
  INSERT INTO urls (short_url, original_url, user_id)
  VALUES
  ($1, $2, $3)
  `

	response := make([]entity.StorageURL, 0, len(batch))

	for _, item := range batch {
		shortURL := random.NewRandomString()

		response = append(response, entity.StorageURL{ID: item.ID, ShortURL: shortURL})

		_, err = tx.Exec(ctx, txQuery, shortURL, item.OriginalURL, userID)
		if err != nil {
			return nil, fmt.Errorf("Storage - StoreBatch - tx.Exec: %w", err)
		}
	}
	tx.Commit(ctx)

	return response, nil
}
