package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lovelydaemon/url-shortener/internal/entity"
)

type Storage struct {
	store  []entity.StorageItem
	file   *os.File
	writer *json.Encoder
	reader *json.Decoder
}

func New(path string) (*Storage, error) {
	storage := &Storage{
		store: make([]entity.StorageItem, 0),
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

	return s.reader.Decode(&s.store)
}

func (s *Storage) Write(item entity.StorageItem) error {
	s.store = append(s.store, item)
	if s.file == nil {
		return nil
	}

	if err := s.file.Truncate(0); err != nil {
		return fmt.Errorf("ShortURLRepo - Store - r.file.Truncate: %w", err)
	}

	if _, err := s.file.Seek(0, 0); err != nil {
		return fmt.Errorf("ShortURLRepo - Store - r.file.Seek: %w", err)
	}

	if err := s.writer.Encode(s.store); err != nil {
		return fmt.Errorf("ShortURLRepo - Store - r.writer.Encode: %w", err)
	}

	return nil
}

func (s *Storage) Get(token string) (entity.StorageItem, bool) {
	for _, v := range s.store {
		if v.Token == token {
			return v, true
		}
	}
	return entity.StorageItem{}, false
}

func (s *Storage) Len() int {
	return len(s.store)
}

func (s *Storage) Close() error {
	return s.file.Close()
}
