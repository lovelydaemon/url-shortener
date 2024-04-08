package v1

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/random"
	urlc "github.com/lovelydaemon/url-shortener/internal/url"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
)

type shortenRoutes struct {
	u usecase.Shorten
	l logger.Interface
}

func NewShortenRoutes(handler *chi.Mux, l logger.Interface, u usecase.Shorten) {
	r := shortenRoutes{u, l}

	handler.Post("/api/shorten", r.generateShortURL)
}

type createShortURLRequest struct {
	URL string `json:"url"`
}

type createShortURLResponse struct {
	Result string `json:"result"`
}

func (r *shortenRoutes) generateShortURL(w http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	contentType := request.Header.Get("Content-Type")
	if contentType != "application/json" {
		r.l.Info("Bad content type", contentType)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var req createShortURLRequest
	dec := json.NewDecoder(request.Body)
	if err := dec.Decode(&req); err != nil {
		r.l.Info("cannot decode request JSON body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := url.ParseRequestURI(req.URL); err != nil {
		r.l.Info("Incorrect body url", req.URL)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := random.NewRandomString(9)

	if err := r.u.Store(ctx, req.URL, token); err != nil {
		r.l.Error(err, "http - v1 - generateShortURL")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := createShortURLResponse{
		Result: urlc.CreateValidURL(request.Host, token),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		r.l.Info("error encoding response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.l.Info("Short url created, 201")
}
