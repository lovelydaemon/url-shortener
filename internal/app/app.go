// Package app configures and runs application
package app

import (
	"fmt"

	"github.com/lovelydaemon/url-shortener/config"
	v1 "github.com/lovelydaemon/url-shortener/internal/controller/http/v1"
	"github.com/lovelydaemon/url-shortener/internal/httpserver"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
	"github.com/lovelydaemon/url-shortener/internal/usecase/repo"
)

// Run creates objects via constructors
func Run(cfg *config.Config) error {

	// Use case
	shortURLUseCase := usecase.New(
		repo.New(),
	)

	// HTTP Server
	r := v1.NewRouter()
	r.Mount("/", v1.NewShortURLRoutes(shortURLUseCase, cfg.ShortAddr))

	httpserver := httpserver.New(r, httpserver.Addr(cfg.Addr))
	fmt.Printf("Server running on %s\n", httpserver.Addr)
	return httpserver.ListenAndServe()
}
