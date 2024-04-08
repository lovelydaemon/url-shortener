// Package app configures and runs application
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/lovelydaemon/url-shortener/config"
	v1 "github.com/lovelydaemon/url-shortener/internal/controller/http/v1"
	"github.com/lovelydaemon/url-shortener/internal/httpserver"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/postgres"
	"github.com/lovelydaemon/url-shortener/internal/storage"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
	"github.com/lovelydaemon/url-shortener/internal/usecase/repo"
)

// Run creates objects via constructors
func Run(cfg *config.Config) {
	// Migrations
	ApplyMigrations(cfg.PG.URL)

	// Logger
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	st, err := storage.New(cfg.Storage.Path)
	if err != nil {
		l.Fatal(fmt.Errorf("Couldn't open file: %w", err))
	}
	defer st.Close()

	var shortenRepo usecase.ShortenRepo
	if cfg.PG.URL != "" {
		shortenRepo = repo.NewShortenPG(pg)
	} else {
		shortenRepo = repo.NewShortenST(st)
	}

	// Use case
	shortenUC := usecase.NewShorten(shortenRepo)
	pingUC := usecase.NewPing(repo.NewPing(pg))

	// HTTP Server
	handler := chi.NewRouter()
	handler = v1.NewRouter(handler, l)
	v1.NewShortURLRoutes(handler, l, shortenUC, cfg.BaseURL)
	v1.NewShortenRoutes(handler, l, shortenUC)
	v1.NewPingRoutes(handler, l, pingUC)

	httpServer := httpserver.New(handler, httpserver.Addr(cfg.HTTP.Addr))
	l.Info("Server running on " + httpServer.Addr())

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
