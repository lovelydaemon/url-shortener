package httpserver

import (
	"net"
	"net/http"
)

type Option func(*http.Server)

func Port(port string) Option {
	return func(s *http.Server) {
		s.Addr = net.JoinHostPort("", port)
	}
}
