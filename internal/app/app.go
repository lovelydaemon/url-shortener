// Package app configures and runs application
package app

import (
	v1 "github.com/lovelydaemon/url-shortener/internal/controller/http/v1"
	"github.com/lovelydaemon/url-shortener/internal/httpserver"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
	"github.com/lovelydaemon/url-shortener/internal/usecase/repo"
)

// Run creates objects via constructors
func Run() error {

	// Use case
	shortURLUseCase := usecase.New(
		repo.New(),
	)

	// HTTP Server
	r := v1.NewRouter(shortURLUseCase)
	httpserver := httpserver.New(r)
	return httpserver.ListenAndServe()
}
