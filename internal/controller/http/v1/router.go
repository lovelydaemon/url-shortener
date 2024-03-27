package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
)

// NewRouter -
func NewRouter(handler *gin.Engine, u usecase.ShortURL) {
	newShortURLRoutes(handler, u)
}
