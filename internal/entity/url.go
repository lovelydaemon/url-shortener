package entity

import (
	"fmt"
	"net/url"
	"strings"
)

type URL string

func NewURL(baseURL string, path string) URL {
	if path == "" {
		return URL(baseURL)
	}

	path = strings.TrimPrefix(path, "/")
	if strings.HasPrefix(baseURL, "http") {
		return URL(fmt.Sprintf("%s/%s", baseURL, path))
	}
	return URL(fmt.Sprintf("http://%s/%s", baseURL, path))
}

func (u URL) Validate() error {
	_, err := url.ParseRequestURI(string(u))
	return err
}
