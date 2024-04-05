package v1

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lovelydaemon/url-shortener/internal/logger"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
)

type pingRoutes struct {
	u usecase.Ping
	l logger.Interface
}

func NewPingRoutes(handler *chi.Mux, u usecase.Ping, l logger.Interface) {
	r := pingRoutes{u, l}

	handler.Get("/ping", r.ping)
}

func (r *pingRoutes) ping(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	if err := r.u.Ping(ctx); err != nil {
		r.l.Info("Couldn't ping database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
