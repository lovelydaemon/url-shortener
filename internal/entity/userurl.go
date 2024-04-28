package entity

import "github.com/google/uuid"

type UserURL struct {
	ShortURL    URL `json:"short_url"`
	OriginalURL URL `json:"original_url"`
}

type UserWithURLs struct {
	UserID    uuid.UUID
	ShortURLs []string
}
