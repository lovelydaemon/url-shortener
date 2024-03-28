package httpserver

import "net/http"

type Option func(*http.Server)

func Addr(addr string) Option {
	return func(s *http.Server) {
		s.Addr = addr
	}
}
