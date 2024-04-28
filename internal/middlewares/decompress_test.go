package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestDecompress(t *testing.T) {
	handler := chi.NewRouter()
	handler.Use(RequestDecompress)

	handler.Post("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})

	srv := httptest.NewServer(handler)
	defer srv.Close()

	originalString := "hello"

	var b bytes.Buffer
	gw := gzip.NewWriter(&b)
	_, err := gw.Write([]byte(originalString))
	require.NoError(t, err, "Error compress data")

	gw.Close()

	body, err := io.ReadAll(&b)
	require.NoError(t, err, "Error read from buffer")

	tests := []struct {
		name            string
		contentType     string
		contentEncoding string
		body            []byte
		expectedCode    int
		expectedBody    string
	}{
		{
			name:            "success decompress",
			contentType:     "application/x-gzip",
			contentEncoding: "gzip",
			body:            body,
			expectedCode:    http.StatusOK,
			expectedBody:    originalString,
		},
		{
			name:         "request without compressing",
			contentType:  "application/json",
			body:         []byte(originalString),
			expectedCode: http.StatusOK,
			expectedBody: originalString,
		},
		{
			name:            "corrupted gzip data",
			contentType:     "application/x-gzip",
			contentEncoding: "gzip",
			body:            []byte(originalString),
			expectedCode:    http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := resty.New().R()
			client.SetBody(tt.body)
			client.SetHeader("Content-Type", tt.contentType)

			if tt.contentEncoding != "" {
				client.SetHeader("Content-Encoding", tt.contentEncoding)
			}

			resp, err := client.Post(srv.URL)
			assert.NoError(t, err, "Error making HTTP request")

			assert.Equal(t, tt.expectedCode, resp.StatusCode(), "Response code didn't match expected")

			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, string(resp.Body()))
			}
		})
	}
}
