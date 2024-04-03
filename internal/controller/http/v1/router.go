package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lovelydaemon/url-shortener/config"
	"github.com/lovelydaemon/url-shortener/internal/controller/http/middlewares"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
)

func NewRouter(u usecase.ShortURL, l logger.Interface, cfg *config.Config) *chi.Mux {
	handler := chi.NewRouter()

	// Middlewares
	handler.Use(middlewares.Logger(l))
	handler.Use(middleware.Recoverer)
	handler.Use(middlewares.RequestDecompress)
	handler.Use(middleware.Compress(5, "application/json", "text/html"))

	NewShortURLRoutes(handler, u, l, cfg.BaseURL)
	NewShortenRoutes(handler, u, l)

	return handler
}
