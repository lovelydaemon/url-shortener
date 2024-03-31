package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lovelydaemon/url-shortener/internal/logger"
)

type basicWriter struct {
	http.ResponseWriter
	bytes int
	code  int
}

func (b *basicWriter) Write(p []byte) (int, error) {
	size, err := b.ResponseWriter.Write(p)
	b.bytes += size
	return size, err
}

func (b *basicWriter) WriteHeader(statusCode int) {
	b.ResponseWriter.WriteHeader(statusCode)
	b.code = statusCode
}

func (b *basicWriter) Status() int {
	return b.code
}

func (b *basicWriter) BytesWritten() int {
	return b.bytes
}

func Logger(l logger.Interface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			bw := &basicWriter{ResponseWriter: w}

			t1 := time.Now()

			defer func() {
				l.Info(
					fmt.Sprintf("%s http://%s%s %s from %s - %d %dB in %s",
						r.Method,
						r.Host,
						r.RequestURI,
						r.Proto,
						r.RemoteAddr,
						bw.Status(),
						bw.BytesWritten(),
						time.Since(t1)))
			}()
			next.ServeHTTP(bw, r)
		}
		return http.HandlerFunc(fn)
	}
}
