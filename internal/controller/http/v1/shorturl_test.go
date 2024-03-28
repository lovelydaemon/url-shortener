package v1

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/lovelydaemon/url-shortener/internal/usecase"
	"github.com/lovelydaemon/url-shortener/internal/usecase/repo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func executeRequest(req *http.Request, r *chi.Mux) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func Test_shortURLRoutes_getOriginalURL(t *testing.T) {
	usecase := usecase.New(repo.New())
	r := chi.NewRouter()
	r.Mount("/", NewShortURLRoutes(usecase, "example.com:8080"))

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://google.com"))
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	res := executeRequest(req, r)

	bodyData, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	url, err := url.ParseRequestURI(string(bodyData))
	require.NoError(t, err)

	token := strings.TrimLeft(url.Path, "/")

	t.Run("valid redirect", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example.com:8080/"+token, nil)
		res := executeRequest(req, r)
		assert.Equal(t, http.StatusTemporaryRedirect, res.Code)
	})
	t.Run("invalid not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/example", nil)

		res := executeRequest(req, r)
		assert.Equal(t, http.StatusNotFound, res.Code)
	})
}

func Test_shortURLRoutes_createShortURL(t *testing.T) {
	usecase := usecase.New(repo.New())
	r := chi.NewRouter()
	r.Mount("/", NewShortURLRoutes(usecase, "localhost:8080"))

	cases := []struct {
		name         string
		method       string
		url          string
		endpoint     string
		contentType  string
		expectedCode int
	}{
		{
			name:         "invalid bad content type",
			method:       http.MethodPost,
			url:          "http://example.com",
			endpoint:     "/",
			contentType:  "application/json",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid bad body data",
			method:       http.MethodPost,
			url:          "example.com",
			endpoint:     "/",
			contentType:  "text/plain; charset=utf-8",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "valid first time created",
			method:       http.MethodPost,
			url:          "https://example.com",
			endpoint:     "/",
			contentType:  "text/plain; charset=utf-8",
			expectedCode: http.StatusCreated,
		},
		{
			name:         "valid url already exists",
			method:       http.MethodPost,
			url:          "https://example.com",
			endpoint:     "/",
			contentType:  "text/plain; charset=utf-8",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.endpoint, strings.NewReader(tt.url))
			req.Header.Set("Content-Type", tt.contentType)

			res := executeRequest(req, r)
			assert.Equal(t, tt.expectedCode, res.Code)
		})
	}
}
