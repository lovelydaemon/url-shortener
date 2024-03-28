package v1

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lovelydaemon/url-shortener/internal/rnd"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
	"github.com/lovelydaemon/url-shortener/internal/validation"
)

type shortURLRoutes struct {
	u usecase.ShortURL
}

func newShortURLRoutes(u usecase.ShortURL) *chi.Mux {
	r := &shortURLRoutes{u}
	router := chi.NewRouter()

	router.Get("/{token}", r.getOriginalURL)
	router.Post("/", r.createShortURL)

	return router
}

func (r *shortURLRoutes) getOriginalURL(w http.ResponseWriter, req *http.Request) {
	token := chi.URLParam(req, "token")
  url := fmt.Sprintf("http://%s/%s", req.Host, token)

	if u, ok := r.u.Get(url); ok {
		w.Header().Set("Location", u)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	http.Error(w, "not found", http.StatusNotFound)
}

func (r *shortURLRoutes) createShortURL(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")

	if contentType != "text/plain; charset=utf-8" {
		http.Error(w, "bad content type", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if validation.IsValidUrl(string(body)) != nil {
		http.Error(w, "bad body data", http.StatusBadRequest)
		return
	}

	if shortURL, ok := r.u.Get(string(body)); ok {
		w.Write([]byte(shortURL))
		return
	}

  shortURL := fmt.Sprintf("http://%s/%s", req.Host, rnd.NewRandomString(9))
	r.u.Create(string(body), shortURL)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}
