package storage

import (
	"fmt"

	"github.com/lovelydaemon/url-shortener/config"
	"github.com/lovelydaemon/url-shortener/internal/pkg/logger"
	"github.com/lovelydaemon/url-shortener/internal/repository"
	"github.com/lovelydaemon/url-shortener/internal/storage/filestorage"
	"github.com/lovelydaemon/url-shortener/internal/storage/inmemory"
	"github.com/lovelydaemon/url-shortener/internal/storage/postgres"
)

type Storage interface {
	repository.Ping
	repository.User
	repository.Shorten
	Close()
}

// New creates a new storage for services
func New(cfg *config.Config, l logger.Interface) Storage {
	if cfg.PG.URL != "" {
		storage, err := postgres.NewStorage(cfg.PG.URL)
		if err == nil {
			l.Info("Used postgres storage")
			return storage
		}
		l.Error(fmt.Errorf("Storage - New - postgres.NewStorage: %w", err))
	}

	if cfg.Storage.Path != "" {
		storage, err := filestorage.NewStorage(cfg.Storage.Path)
		if err == nil {
			l.Info("Used file storage")
			return storage
		}
		l.Error(fmt.Errorf("Storage - New - filestorage.NewStorage: %w", err))
	}

	l.Info("Used in-memory storage")
	return inmemory.NewStorage()
}
