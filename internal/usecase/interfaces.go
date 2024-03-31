package usecase

type (
	ShortURL interface {
		Get(url string) (string, bool)
		Create(originalUrl, token string)
	}

	ShortURLRepo interface {
		Get(url string) (string, bool)
		Create(originalUrl, token string)
	}
)
