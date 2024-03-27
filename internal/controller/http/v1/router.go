package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
)

func NewRouter(u usecase.ShortURL) *chi.Mux {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// shortURL get and create
	r.Mount("/", newShortURLRoutes(u))
	return r
}
