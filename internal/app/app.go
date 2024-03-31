// Package app configures and runs application
package app

import (
	"github.com/lovelydaemon/url-shortener/config"
	v1 "github.com/lovelydaemon/url-shortener/internal/controller/http/v1"
	"github.com/lovelydaemon/url-shortener/internal/httpserver"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
	"github.com/lovelydaemon/url-shortener/internal/usecase/repo"
)

// Run creates objects via constructors
func Run(cfg *config.Config) error {
	l := logger.New(cfg.Log.Level)

	// Use case
	shortURLUseCase := usecase.New(
		repo.New(),
	)

	// HTTP Server
	handler := v1.NewRouter(shortURLUseCase, l, cfg)
	httpServer := httpserver.New(handler, httpserver.Addr(cfg.HTTP.Addr))
	l.Info("Server running on " + httpServer.Addr)
	return httpServer.ListenAndServe()
}
