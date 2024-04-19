package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/url"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
)

type userRoutes struct {
	u usecase.User
	l logger.Interface
}

func NewUserRoutes(handler chi.Router, l logger.Interface, u usecase.User) {
	r := userRoutes{u, l}

	handler.Get("/api/user/urls", r.getUserURLs)
	handler.Delete("/api/user/urls", r.deleteURLs)
}

func (r *userRoutes) getUserURLs(w http.ResponseWriter, request *http.Request) {
	if _, ok := request.Context().Value("userID").(uuid.UUID); !ok {
		r.l.Info("Empty user id")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp, err := r.u.GetURLs(request.Context())
	if err != nil {
		r.l.Error(fmt.Errorf("Error get user urls: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(resp) == 0 {
		r.l.Info("Response data is empty")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	for i, v := range resp {
		v.ShortURL = url.CreateValidURL(request.Host, v.ShortURL)
		resp[i] = v
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		r.l.Info("error encoding response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (r *userRoutes) deleteURLs(w http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	if _, ok := ctx.Value("userID").(uuid.UUID); !ok {
		r.l.Info("Empty user id")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req []string
	dec := json.NewDecoder(request.Body)
	if err := dec.Decode(&req); err != nil {
		r.l.Info("cannot decode request JSON body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(req) == 0 {
		r.l.Info("request body is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	r.u.DeleteURLs(ctx, req)

	w.WriteHeader(http.StatusAccepted)
}
