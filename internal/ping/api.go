package ping

import (
	"github.com/go-chi/chi/v5"
	"github.com/lovelydaemon/url-shortener/internal/pkg/logger"
	"net/http"
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

	r.Get("/ping", res.ping)
}

func (r resource) ping(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if err := r.service.Ping(ctx); err != nil {
		r.log.Info("Could not ping database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
