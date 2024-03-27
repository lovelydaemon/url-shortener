package v1

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
	"github.com/lovelydaemon/url-shortener/internal/usecase/repo"
	"github.com/stretchr/testify/assert"
)

func Test_shortURLRoutes_createShortURL(t *testing.T) {
  gin.SetMode("test")
	r := gin.Default()
	newShortURLRoutes(r, usecase.New(repo.New()))

	t.Run("invalid bad content type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://example.com"))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("invalid bad url link", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("example.com"))
		req.Header.Set("Content-Type", "text/plain")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("valid first time create", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://example.com"))
		req.Header.Set("Content-Type", "text/plain")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
	t.Run("valid already exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("http://example.com"))
		req.Header.Set("Content-Type", "text/plain")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
