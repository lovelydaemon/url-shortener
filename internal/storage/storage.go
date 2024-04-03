package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lovelydaemon/url-shortener/internal/entity"
)

type Storage struct {
	Store  []entity.StorageItem
	File   *os.File
	Writer *json.Encoder
	Reader *json.Decoder
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

	storage.File = file
	storage.Writer = json.NewEncoder(file)
	storage.Reader = json.NewDecoder(file)
	storage.Reader.UseNumber()

	if err := storage.LoadStore(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) LoadStore() error {
	fInfo, err := s.File.Stat()
	if err != nil {
		return err
	}

	if fInfo.Size() == 0 {
		return nil
	}

	return s.Reader.Decode(&s.Store)
}

func (s *Storage) Write(item entity.StorageItem) error {
	s.Store = append(s.Store, item)
	if s.File == nil {
		return nil
	}

	if err := s.File.Truncate(0); err != nil {
		return fmt.Errorf("ShortURLRepo - Store - r.file.Truncate: %w", err)
	}

	if _, err := s.File.Seek(0, 0); err != nil {
		return fmt.Errorf("ShortURLRepo - Store - r.file.Seek: %w", err)
	}

	if err := s.Writer.Encode(s.Store); err != nil {
		return fmt.Errorf("ShortURLRepo - Store - r.writer.Encode: %w", err)
	}

	return nil
}

func (s *Storage) Close() error {
	return s.File.Close()
}
