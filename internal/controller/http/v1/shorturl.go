package v1

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lovelydaemon/url-shortener/internal/rnd"
	"github.com/lovelydaemon/url-shortener/internal/url"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
	"github.com/lovelydaemon/url-shortener/internal/validation"
)

type shortURLRoutes struct {
	u         usecase.ShortURL
	shortAddr string
}

func NewShortURLRoutes(u usecase.ShortURL, shortAddr string) *chi.Mux {
	r := &shortURLRoutes{u, shortAddr}
	router := chi.NewRouter()

	router.Get("/{token}", r.getOriginalURL)
	router.Post("/", r.createShortURL)

	return router
}

func (r *shortURLRoutes) getOriginalURL(w http.ResponseWriter, req *http.Request) {
	token := chi.URLParam(req, "token")

	if u, ok := r.u.Get(token); ok {
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

	if token, ok := r.u.Get(string(body)); ok {
    shortURL := url.CreateValidURL(r.shortAddr, token)
		w.Write([]byte(shortURL))
		return
	}
  
	token := rnd.NewRandomString(9)
	r.u.Create(string(body), token)
  
  shortURL := url.CreateValidURL(r.shortAddr, token)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}
