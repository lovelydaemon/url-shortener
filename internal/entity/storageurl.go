package entity

import (
	"time"

	"github.com/google/uuid"
)

type StorageURL struct {
	ID          uuid.UUID
	ShortURL    string
	OriginalURL URL
	UserID      uuid.UUID
	CreatedAt   time.Time
	DeletedFlag bool
}
