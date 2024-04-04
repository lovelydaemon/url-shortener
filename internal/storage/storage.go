package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lovelydaemon/url-shortener/internal/entity"
)

type Storage struct {
	Store  []entity.StorageItem
	file   *os.File
	writer *json.Encoder
	reader *json.Decoder
}

func NewStorage(path string) (*Storage, error) {
	storage := &Storage{
		Store: make([]entity.StorageItem, 0),
	}
	if path == "" {
		return storage, nil
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	storage.file = file
	storage.writer = json.NewEncoder(file)
	storage.reader = json.NewDecoder(file)
	storage.reader.UseNumber()

	if err := storage.LoadStore(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) LoadStore() error {
	fInfo, err := s.file.Stat()
	if err != nil {
		return err
	}

	if fInfo.Size() == 0 {
		return nil
	}

	return s.reader.Decode(&s.Store)
}

func (s *Storage) Write(item entity.StorageItem) error {
	s.Store = append(s.Store, item)
	if s.file == nil {
		return nil
	}

	if err := s.file.Truncate(0); err != nil {
		return fmt.Errorf("ShortURLRepo - Store - r.file.Truncate: %w", err)
	}

	if _, err := s.file.Seek(0, 0); err != nil {
		return fmt.Errorf("ShortURLRepo - Store - r.file.Seek: %w", err)
	}

	if err := s.writer.Encode(s.Store); err != nil {
		return fmt.Errorf("ShortURLRepo - Store - r.writer.Encode: %w", err)
	}

	return nil
}

func (s *Storage) Close() error {
	return s.file.Close()
}
