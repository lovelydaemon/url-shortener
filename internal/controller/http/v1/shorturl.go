package v1

import (
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/rnd"
	urlc "github.com/lovelydaemon/url-shortener/internal/url"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
)

type shortURLRoutes struct {
	u         usecase.ShortURL
	l         logger.Interface
	shortAddr string
}

func NewShortURLRoutes(handler *chi.Mux, u usecase.ShortURL, l logger.Interface, shortAddr string) {
	r := &shortURLRoutes{u, l, shortAddr}

	handler.Get("/{token}", r.getOriginalURL)
	handler.Post("/", r.createShortURL)
}

func (r *shortURLRoutes) getOriginalURL(w http.ResponseWriter, req *http.Request) {
	token := chi.URLParam(req, "token")

	if u, ok := r.u.Get(token); ok {
		r.l.Info("Found original url", u)
		w.Header().Set("Location", u)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	r.l.Info("Original url not found")
	w.WriteHeader(http.StatusNotFound)
}

func (r *shortURLRoutes) createShortURL(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		r.l.Info("Error reading body", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bodyURL := string(body)

	if _, err := url.ParseRequestURI(bodyURL); err != nil {
		r.l.Info("Invalid request body", bodyURL)
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	token := rnd.NewRandomString(9)
	r.u.Create(bodyURL, token)

	var baseURL string
	if r.shortAddr != "" {
		baseURL = r.shortAddr
	} else {
		baseURL = req.Host
	}

	shortURL := urlc.CreateValidURL(baseURL, token)

	r.l.Info("Short url created, 201")
  w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}
