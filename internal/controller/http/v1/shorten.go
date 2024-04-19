package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lovelydaemon/url-shortener/internal/entity"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/url"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
)

type shortenRoutes struct {
	u usecase.Shorten
	l logger.Interface
}

func NewShortenRoutes(handler chi.Router, l logger.Interface, u usecase.Shorten) {
	r := shortenRoutes{u, l}

	handler.Route("/api/shorten", func(h chi.Router) {
		h.Use(middleware.AllowContentType("application/json"))
		h.Post("/", r.generateShortURL)
		h.Post("/batch", r.generateShortURLBatch)
	})
}

type generateShortURLRequest struct {
	URL string `json:"url"`
}

type generateShortURLResponse struct {
	Result string `json:"result"`
}

func (r *shortenRoutes) generateShortURL(w http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	var req generateShortURLRequest
	dec := json.NewDecoder(request.Body)
	if err := dec.Decode(&req); err != nil {
		r.l.Info("cannot decode request JSON body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := url.Validate(req.URL); err != nil {
		r.l.Info("bad original url", req.URL)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL, err := r.u.Store(ctx, req.URL)
	if err != nil && !errors.Is(err, ErrConflict) {
		r.l.Error(err, "http - v1 - generateShortURL")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := generateShortURLResponse{
		Result: url.CreateValidURL(request.Host, shortURL),
	}

	w.Header().Set("Content-Type", "application/json")

	if errors.Is(err, ErrConflict) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		r.l.Info("error encoding response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (r *shortenRoutes) generateShortURLBatch(w http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	var req []entity.BatchItemIn
	dec := json.NewDecoder(request.Body)
	if err := dec.Decode(&req); err != nil {
		r.l.Info("cannot decode request JSON body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(req) == 0 {
		r.l.Info("empty request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, v := range req {
		if err := url.Validate(v.OriginalURL); err != nil {
			r.l.Info(fmt.Sprintf("bad original url %s, id %s", v.OriginalURL, v.ID))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	resp, err := r.u.StoreBatch(ctx, req)
	if err != nil {
		r.l.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for i, v := range resp {
		v.ShortURL = url.CreateValidURL(request.Host, v.ShortURL)
		resp[i] = v
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		r.l.Info("error encoding response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
