package v1

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/url"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
)

type shortURLRoutes struct {
	u         usecase.Shorten
	l         logger.Interface
	shortAddr string
}

func NewShortURLRoutes(handler *chi.Mux, l logger.Interface, u usecase.Shorten, shortAddr string) {
	r := &shortURLRoutes{u, l, shortAddr}

	handler.Get("/{token}", r.getOriginalURL)
	handler.Post("/", r.generateShortURL)
}

func (r *shortURLRoutes) getOriginalURL(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	token := chi.URLParam(req, "token")

	item, err := r.u.Get(ctx, token)
	if err != nil {
		r.l.Info("Error get url by token", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	r.l.Info("Found original url", item.OriginalURL)
	w.Header().Set("Location", item.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	return
}

func (r *shortURLRoutes) generateShortURL(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		r.l.Info("Error reading body", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	originalURL := string(body)
	if err := url.Validate(originalURL); err != nil {
		r.l.Info("Invalid request body", originalURL)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := r.u.Store(ctx, originalURL)
	if err != nil {
		r.l.Error(err, "http - v1 - generateShortURL")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var baseURL string
	if r.shortAddr != "" {
		baseURL = r.shortAddr
	} else {
		baseURL = req.Host
	}

	shortURL := url.CreateValidURL(baseURL, token)

	r.l.Info("Short url created, 201")
	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}
