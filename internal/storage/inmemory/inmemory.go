package inmemory

import (
	"sync"

	"github.com/lovelydaemon/url-shortener/internal/entity"
)

type storage struct {
	mu      sync.RWMutex
	storage map[string]entity.StorageURL
}

// NewStorage creates new in-memory storage
func NewStorage() *storage {
	storage := &storage{
		mu:      sync.RWMutex{},
		storage: make(map[string]entity.StorageURL),
	}
	return storage
}

func (s *storage) Close() {}
