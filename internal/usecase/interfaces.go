package usecase

type (
	ShortURL interface {
		Get(url string) (string, bool)
		Create(originalUrl, shortURL string)
	}

	ShortURLRepo interface {
		Get(url string) (string, bool)
		Create(originalUrl, shortURL string)
	}
)
