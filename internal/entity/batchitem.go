package entity

type BatchItemIn struct {
	ID          string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type BatchItemOut struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}
