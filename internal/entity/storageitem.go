// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

type StorageItem struct {
	ID          int    `json:"id"`
	Token       string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
