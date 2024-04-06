// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

type StorageItem struct {
	UUID        int    `json:"uuid"`
	Token       string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
