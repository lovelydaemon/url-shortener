package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lovelydaemon/url-shortener/config"
	"github.com/lovelydaemon/url-shortener/internal/middlewares"
	"github.com/lovelydaemon/url-shortener/internal/ping"
	"github.com/lovelydaemon/url-shortener/internal/pkg/logger"
	"github.com/lovelydaemon/url-shortener/internal/queue"
	"github.com/lovelydaemon/url-shortener/internal/shorten"
	"github.com/lovelydaemon/url-shortener/internal/storage"
	"github.com/lovelydaemon/url-shortener/internal/user"
)

func BuildHandler(l logger.Interface, cfg *config.Config, s storage.Storage, q *queue.Queue) *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		middlewares.Logger(l),
		middleware.Recoverer,
		middlewares.RequestDecompress,
		middleware.Compress(5, "application/json", "text/html"),
	)

	// Public
	router.Group(func(r chi.Router) {
		r.Use(middlewares.Authentication(cfg.JWT.Key))
		shorten.RegisterHandlers(r,
			shorten.NewService(s),
			l,
			cfg.BaseURL,
		)
		ping.RegisterHandlers(r,
			ping.NewService(s),
			l,
		)
	})

	// Private
	router.Group(func(r chi.Router) {
		r.Use(middlewares.Authorization(cfg.JWT.Key, l))
		user.RegisterHandlers(r,
			user.NewService(s, q),
			l,
		)
	})

	return router
}
