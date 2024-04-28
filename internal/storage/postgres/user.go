package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lovelydaemon/url-shortener/internal/entity"
)

func (s storage) GetURLs(ctx context.Context) ([]entity.UserURL, error) {
	userID := ctx.Value("userID").(uuid.UUID)

	query := `
  SELECT short_url, original_url FROM urls
  WHERE user_id = $1
  `

	rows, err := s.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("Storage - GetURLs - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var response []entity.UserURL

	for rows.Next() {
		var userURL entity.UserURL
		if err := rows.Scan(&userURL.ShortURL, &userURL.OriginalURL); err != nil {
			return nil, fmt.Errorf("Storage - GetURLs - rows.Scan: %w", err)
		}

		response = append(response, userURL)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("Storage - GetURLs - rows.Err: %w", err)
	}

	return response, nil
}

func (s storage) DeleteURLs(ctx context.Context, data ...entity.UserWithURLs) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Storage - DeleteURLs - s.Pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
  UPDATE urls
  SET is_deleted = TRUE
  WHERE short_url = ANY($1) AND user_id = $2
  `

	for _, v := range data {
		if _, err := tx.Exec(ctx, query, v.ShortURLs, v.UserID); err != nil {
			return fmt.Errorf("Storage - DeleteURLs - tx.Exec: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("Storage - DeleteURLs - tx.Commit: %w", err)
	}
	return nil
}
