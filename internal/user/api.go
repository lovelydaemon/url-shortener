package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lovelydaemon/url-shortener/internal/entity"
	"github.com/lovelydaemon/url-shortener/internal/pkg/logger"
)

type resource struct {
	service Service
	log     logger.Interface
}

func RegisterHandlers(r chi.Router, service Service, l logger.Interface) {
	res := resource{
		service: service,
		log:     l,
	}

	r.Get("/api/user/urls", res.getUserURLs)
	r.Delete("/api/user/urls", res.deleteUserURLs)
}

func (r resource) getUserURLs(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	response, err := r.service.GetURLs(ctx)
	if err != nil {
		r.log.Error(fmt.Errorf("Error get user urls: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(response) == 0 {
		r.log.Info("Response data is empty")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	for i, item := range response {
		item.ShortURL = entity.NewURL(req.Host, string(item.ShortURL))
		response[i] = item
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	if err := enc.Encode(response); err != nil {
		r.log.Info("Error encoding response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (r resource) deleteUserURLs(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var request []string
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&request); err != nil {
		r.log.Info("Cannot decode request JSON body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(request) > 0 {
		r.service.DeleteURLs(ctx, request)
	}

	w.WriteHeader(http.StatusAccepted)
}
