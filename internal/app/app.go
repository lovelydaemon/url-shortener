// Package app configures and runs application
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lovelydaemon/url-shortener/config"
	"github.com/lovelydaemon/url-shortener/internal/controller/http"
	"github.com/lovelydaemon/url-shortener/internal/pkg/httpserver"
	"github.com/lovelydaemon/url-shortener/internal/pkg/logger"
	"github.com/lovelydaemon/url-shortener/internal/queue"
	"github.com/lovelydaemon/url-shortener/internal/storage"
)

// Run creates objects via constructors
func Run(cfg *config.Config) {
	// Migrations
	ApplyMigrations(cfg.PG.URL)

	// Logger
	l := logger.New(cfg.Log.Level)

	// Storage
	storage := storage.New(cfg, l)
	defer storage.Close()

	// Queue
	queue := queue.New(storage, l)
	ctx, cancel := context.WithCancel(context.Background())
	go queue.FlushUserURLs(ctx)

	// HTTP Server
	handler := http.BuildHandler(l, cfg, storage, queue)
	httpServer := httpserver.New(handler, httpserver.Addr(cfg.HTTP.Addr))
	l.Info("Server running on " + httpServer.Addr())

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	cancel()

	if err := httpServer.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
