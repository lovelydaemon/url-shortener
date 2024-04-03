package entity

type StorageItem struct {
	UUID        int    `json:"uuid"`
	Token       string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
