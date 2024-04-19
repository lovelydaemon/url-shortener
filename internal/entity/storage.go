// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import "github.com/google/uuid"

type Storage struct {
	ID          int    `json:"id" db:"id"`
	ShortURL    string `json:"short_url" db:"short_url"`
	OriginalURL string `json:"original_url" db:"original_url"`
	DeletedFlag bool   `json:"is_deleted" db:"is_deleted"`
}

type StorageWithUser struct {
	UserID   uuid.UUID
	ShortURL string
}
