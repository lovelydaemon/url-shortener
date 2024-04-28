package shorten

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lovelydaemon/url-shortener/internal/entity"
	e "github.com/lovelydaemon/url-shortener/internal/pkg/errors"
	"github.com/lovelydaemon/url-shortener/internal/pkg/logger"
)

type resource struct {
	service Service
	log     logger.Interface
	baseURL string
}

func RegisterHandlers(r chi.Router, service Service, l logger.Interface, baseURL string) {
	res := resource{
		service: service,
		log:     l,
		baseURL: baseURL,
	}

	r.Get("/{shortURL}", res.getOriginalURL)
	r.Post("/", res.generateShortURL)

	r.Route("/api/shorten", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/", res.generateShortURLJson)
		r.Post("/batch", res.generateShortURLBatch)
	})
}

func (r resource) getOriginalURL(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	shortURL := chi.URLParam(req, "shortURL")

	storageURL, err := r.service.Get(ctx, shortURL)
	if err != nil {
		if errors.Is(err, e.ErrRecNotFound) {
			r.log.Info("Record not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		r.log.Info("Error get original URL", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if storageURL.DeletedFlag {
		r.log.Info("URL was deleted")
		w.WriteHeader(http.StatusGone)
		return
	}

	r.log.Info("Found original URL", storageURL.OriginalURL)
	w.Header().Set("Location", string(storageURL.OriginalURL))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (r resource) generateShortURL(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		r.log.Info("Error reading request body", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	originalURL := entity.NewURL(string(body), "")
	if err := originalURL.Validate(); err != nil {
		r.log.Info("Invalid request body", originalURL)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL, err := r.service.Store(ctx, originalURL)
	if err != nil && !errors.Is(err, e.ErrConflict) {
		r.log.Error(fmt.Errorf("Error generate shortURL: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	baseURL := req.Host
	if r.baseURL != "" {
		baseURL = r.baseURL
	}

	response := entity.NewURL(baseURL, shortURL)
	w.Header().Set("Content-Type", "text/plain")

	if errors.Is(err, e.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte(response))
}

type generateShortURLRequest struct {
	URL entity.URL `json:"url"`
}

type generateShortURLResponse struct {
	Result entity.URL `json:"result"`
}

func (r resource) generateShortURLJson(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var request generateShortURLRequest
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&request); err != nil {
		r.log.Info("Cannot decode request JSON body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := request.URL.Validate(); err != nil {
		r.log.Info("Bad request body", request.URL)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL, err := r.service.Store(ctx, request.URL)
	if err != nil && !errors.Is(err, e.ErrConflict) {
		r.log.Error(fmt.Errorf("Error generate shortURL: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := generateShortURLResponse{
		Result: entity.NewURL(req.Host, shortURL),
	}

	w.Header().Set("Content-Type", "application/json")
	if errors.Is(err, e.ErrConflict) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(response); err != nil {
		r.log.Info("Error encoding response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (r resource) generateShortURLBatch(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var request []BatchRequest
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&request); err != nil {
		r.log.Info("Error decoding request JSON body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(request) == 0 {
		r.log.Info("Empty request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, item := range request {
		if err := item.OriginalURL.Validate(); err != nil {
			r.log.Info("Validation error", item.ID, item.OriginalURL)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	response, err := r.service.StoreBatch(ctx, request)
	if err != nil {
		r.log.Error(fmt.Errorf("Error generate batch response: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for i, item := range response {
		item.ShortURL = entity.NewURL(req.Host, string(item.ShortURL))
		response[i] = item
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	if err := enc.Encode(response); err != nil {
		r.log.Info("Error encoding response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
