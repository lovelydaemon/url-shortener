package v1

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/rnd"
	urlc "github.com/lovelydaemon/url-shortener/internal/url"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
)

type shortenRoutes struct {
	u usecase.ShortURL
	l logger.Interface
}

func NewShortenRoutes(u usecase.ShortURL, l logger.Interface) *chi.Mux {
	r := shortenRoutes{u, l}
	router := chi.NewRouter()

	router.Post("/shorten", r.createShortURL)
	return router
}

type createShortURLRequest struct {
	URL string `json:"url"`
}

type createShortURLResponse struct {
	Result string `json:"result"`
}

func (sr *shortenRoutes) createShortURL(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		sr.l.Info("Bad content type", contentType)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var req createShortURLRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		sr.l.Info("cannot decode request JSON body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := url.ParseRequestURI(req.URL); err != nil {
		sr.l.Info("Incorrect body url", req.URL)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := rnd.NewRandomString(9)
	sr.u.Create(req.URL, token)

	resp := createShortURLResponse{
		Result: urlc.CreateValidURL(r.Host, token),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		sr.l.Info("error encoding response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sr.l.Info("Short url created, 201")
}
