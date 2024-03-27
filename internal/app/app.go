// Package app configures and runs application
package app

import (
	"github.com/gin-gonic/gin"
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
	handler := gin.Default()
	v1.NewRouter(handler, shortURLUseCase)
	httpserver := httpserver.New(handler)
	return httpserver.ListenAndServe()
}
