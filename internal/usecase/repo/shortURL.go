package repo

type ShortURLRepo struct {
	store map[string]string
}

func New() *ShortURLRepo {
	return &ShortURLRepo{
		store: make(map[string]string),
	}
}

func (r *ShortURLRepo) Get(url string) (string, bool) {
	u, ok := r.store[url]
	return u, ok
}

func (r *ShortURLRepo) Create(originalURL, token string) {
	r.store[token] = originalURL
}
