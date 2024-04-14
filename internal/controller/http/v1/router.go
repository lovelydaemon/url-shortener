package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lovelydaemon/url-shortener/config"
	"github.com/lovelydaemon/url-shortener/internal/controller/http/middlewares"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
)

func NewRouter(handler *chi.Mux, l logger.Interface, cfg *config.Config, u *usecase.UseCases) *chi.Mux {
	// Middlewares
	handler.Use(middlewares.Logger(l))
	handler.Use(middleware.Recoverer)
	handler.Use(middlewares.RequestDecompress)
	handler.Use(middleware.Compress(5, "application/json", "text/html"))

	// Public
	handler.Group(func(handler chi.Router) {
		handler.Use(middlewares.Authentication(cfg.JWT.Key))
		NewShortURLRoutes(handler, l, u.Shorten, cfg.BaseURL)
		NewShortenRoutes(handler, l, u.Shorten)
		NewPingRoutes(handler, l, u.Ping)
	})

	// Private
	// Require Authentication
	handler.Group(func(handler chi.Router) {
		handler.Use(middlewares.Authorization(cfg.JWT.Key))
		NewUserRoutes(handler, l, u.User)
	})

	return handler
}
