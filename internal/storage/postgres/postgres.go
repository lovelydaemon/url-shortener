package postgres

import (
	"github.com/lovelydaemon/url-shortener/internal/pkg/postgres"
)

type storage struct {
	*postgres.Postgres
}

// NewStorage creates new postgres storage
func NewStorage(url string) (*storage, error) {
	pg, err := postgres.New(url)
	if err != nil {
		return nil, err
	}

	s := &storage{pg}
	return s, nil
}
