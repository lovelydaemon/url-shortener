package v1

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lovelydaemon/url-shortener/internal/rnd"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
	"github.com/lovelydaemon/url-shortener/internal/validation"
)

type shortURLRoutes struct {
	u usecase.ShortURL
}

func newShortURLRoutes(handler *gin.Engine, u usecase.ShortURL) {
	r := &shortURLRoutes{u}

	handler.GET("/:token", r.getOriginalURL)
	handler.POST("/", r.createShortURL)
}

func (r *shortURLRoutes) getOriginalURL(c *gin.Context) {
	token := c.Param("token")
	url := fmt.Sprintf("%s/%s", c.Request.Host, token)

	u, ok := r.u.Get(url)
	if !ok {
		c.String(http.StatusBadRequest, "not found")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, u)
}

func (r *shortURLRoutes) createShortURL(c *gin.Context) {
	if c.ContentType() != "text/plain" {
		c.String(http.StatusBadRequest, "bad content type")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if validation.IsValidUrl(string(body)) != nil {
		c.String(http.StatusBadRequest, "bad url link")
		return
	}

	shortURL, ok := r.u.Get(string(body))
	if ok {
		c.String(http.StatusOK, shortURL)
		return
	}

	shortURL = fmt.Sprintf("%s/%s", c.Request.Host, rnd.NewRandomString(9))
	r.u.Create(string(body), shortURL)
	c.String(http.StatusCreated, shortURL)
}
