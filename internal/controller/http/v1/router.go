package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lovelydaemon/url-shortener/internal/controller/http/middlewares"
	"github.com/lovelydaemon/url-shortener/internal/logger"
)

func NewRouter(handler *chi.Mux, l logger.Interface) *chi.Mux {
	// Middlewares
	handler.Use(middlewares.Logger(l))
	handler.Use(middleware.Recoverer)
	handler.Use(middlewares.RequestDecompress)
	handler.Use(middleware.Compress(5, "application/json", "text/html"))

	return handler
}
