package filestorage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/lovelydaemon/url-shortener/internal/entity"
)

var (
	ErrFileFormat = errors.New("Unsupported file format")
)

type storage struct {
	mu       sync.RWMutex
	storage  map[string]entity.StorageURL
	filename string
}

// NewStorage creates new file storage
func NewStorage(path string) (*storage, error) {
	filename := filepath.Base(path)

	if filepath.Ext(filename) != ".json" || len(filename) < 6 {
		return nil, ErrFileFormat
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	storage := &storage{
		mu:       sync.RWMutex{},
		storage:  make(map[string]entity.StorageURL),
		filename: path,
	}

	if err := storage.ReadFromFile(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *storage) WriteToFile() error {
	data, err := json.Marshal(s.storage)
	if err != nil {
		return err
	}

	f, err := os.Create(s.filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

func (s *storage) ReadFromFile() error {
	f, err := os.Open(s.filename)
	if err != nil {
		return nil
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	if len(data) != 0 {
		return json.Unmarshal(data, &s.storage)
	}
	return nil
}

func (s *storage) Close() {}
