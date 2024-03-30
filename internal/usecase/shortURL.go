package usecase

type ShortURLUseCase struct {
	repo ShortURLRepo
}

func New(r ShortURLRepo) *ShortURLUseCase {
	return &ShortURLUseCase{
		repo: r,
	}
}

func (uc *ShortURLUseCase) Get(url string) (string, bool) {
	u, ok := uc.repo.Get(url)
	return u, ok
}

func (uc *ShortURLUseCase) Create(originalURL, token string) {
	uc.repo.Create(originalURL, token)
}
