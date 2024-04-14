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

func (r *UserRepo) GetUrls(ctx context.Context) ([]entity.UserURL, error) {
	userID := ctx.Value("userID")

	rows, err := r.Pool.Query(ctx, `
    SELECT short_url, original_url from urls
    WHERE user_id = $1
  `, userID)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetUrls - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var userURLs []entity.UserURL
	for rows.Next() {
		var v entity.UserURL
		if err := rows.Scan(&v.ShortURL, &v.OriginalURL); err != nil {
			return nil, fmt.Errorf("UserRepo - GetUrls - rows.Scan: %w", err)
		}

		userURLs = append(userURLs, v)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetUrls - rows.Err: %w", err)
	}

	return userURLs, nil
}
