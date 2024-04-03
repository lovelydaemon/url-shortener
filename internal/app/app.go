// Package app configures and runs application
package app

import (
	"fmt"

	"github.com/lovelydaemon/url-shortener/config"
	v1 "github.com/lovelydaemon/url-shortener/internal/controller/http/v1"
	"github.com/lovelydaemon/url-shortener/internal/httpserver"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/storage"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
	"github.com/lovelydaemon/url-shortener/internal/usecase/repo"
)

// Run creates objects via constructors
func Run(cfg *config.Config) error {
	// Logger
	l := logger.New(cfg.Log.Level)

	// File Storage
	st, err := storage.NewStorage(cfg.Storage.Path)
	if err != nil {
		l.Fatal(fmt.Errorf("Couldn't open file: %w", err))
	}
	defer st.Close()

	// Use case
	shortURLUseCase := usecase.New(repo.New(st))

	// HTTP Server
	handler := v1.NewRouter(shortURLUseCase, l, cfg)
	httpServer := httpserver.New(handler, httpserver.Addr(cfg.HTTP.Addr))
	l.Info("Server running on " + httpServer.Addr)
	return httpServer.ListenAndServe()
}
