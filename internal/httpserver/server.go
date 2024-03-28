package httpserver

import (
	"net/http"
	"time"
)

const (
	_defaultAddr         = ":8080"
	_defaultReadTimeout  = 10 * time.Second
	_defaultWriteTimeout = 10 * time.Second
)

// New creates default http server
func New(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         _defaultAddr,
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
	}
}
