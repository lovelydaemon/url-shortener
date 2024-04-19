package repo

import (
	"context"
	"fmt"

	"github.com/lovelydaemon/url-shortener/internal/entity"
	"github.com/lovelydaemon/url-shortener/internal/postgres"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUser(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) GetURLs(ctx context.Context) ([]entity.UserURL, error) {
	userID := ctx.Value("userID")

	rows, err := r.Pool.Query(ctx, `
    SELECT short_url, original_url from urls
    WHERE user_id = $1
  `, userID)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetURLs - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var userURLs []entity.UserURL
	for rows.Next() {
		var v entity.UserURL
		if err := rows.Scan(&v.ShortURL, &v.OriginalURL); err != nil {
			return nil, fmt.Errorf("UserRepo - GetURLs - rows.Scan: %w", err)
		}

		userURLs = append(userURLs, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetURLs - rows.Err: %w", err)
	}

	return userURLs, nil
}

func (r *UserRepo) DeleteURLs(ctx context.Context, urls ...entity.StorageWithUser) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("UserRepo - DeleteURLs - r.Pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx)

	userMap := make(map[string][]string, len(urls))

	for _, v := range urls {
		urlsList := userMap[v.UserID.String()]
		urlsList = append(urlsList, v.ShortURL)
		userMap[v.UserID.String()] = urlsList
	}

	query := `UPDATE urls SET is_deleted = TRUE WHERE short_url = ANY($1) AND user_id = $2`

	for userID, urls := range userMap {
		if _, err := tx.Exec(ctx, query, urls, userID); err != nil {
			return fmt.Errorf("UserRepo - DeleteURLs - tx.Exec: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("UserRepo - DeleteURLs - tx.Commit: %w", err)
	}

	return nil
}
