package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var allowedGzipContentTypes = []string{
	"application/x-gzip",
	"application/gzip",
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	cr := &compressReader{
		r:  r,
		zr: zr,
	}

	return cr, nil
}

func (c *compressReader) Read(p []byte) (int, error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func RequestDecompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentEncoding := r.Header.Get("Content-Encoding")
		contentType := r.Header.Get("Content-Type")

		if !sendsGzip(contentEncoding, contentType) {
			next.ServeHTTP(w, r)
			return
		}

		cr, err := newCompressReader(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer cr.Close()

		r.Body = cr

		next.ServeHTTP(w, r)
	})
}

func sendsGzip(contentEncoding, contentType string) bool {
	if !strings.Contains(contentEncoding, "gzip") {
		return false
	}
	for _, v := range allowedGzipContentTypes {
		if strings.Contains(contentType, v) {
			return true
		}
	}
	return false
}
