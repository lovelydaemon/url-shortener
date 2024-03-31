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
	shortAddr string
	l         logger.Interface
}

func NewShortURLRoutes(u usecase.ShortURL, shortAddr string, l logger.Interface) *chi.Mux {
	r := &shortURLRoutes{u, shortAddr, l}
	router := chi.NewRouter()

	router.Get("/{token}", r.getOriginalURL)
	router.Post("/", r.createShortURL)

	return router
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
	http.Error(w, "not found", http.StatusNotFound)
}

func (r *shortURLRoutes) createShortURL(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")

	if contentType != "text/plain; charset=utf-8" {
		r.l.Info("Bad content type", contentType)
		http.Error(w, "bad content type", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bodyURL := string(body)

	if _, err := url.ParseRequestURI(bodyURL); err != nil {
		r.l.Info("Incorrect body url", bodyURL)
		http.Error(w, "bad body data", http.StatusBadRequest)
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
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}
