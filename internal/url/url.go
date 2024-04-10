package url

import (
	"fmt"
	"net/url"
	"strings"
)

func Validate(URL string) error {
	_, err := url.ParseRequestURI(URL)
	return err
}

func CreateValidURL(baseURL string, path string) string {
	if strings.HasPrefix(baseURL, "http") {
		return fmt.Sprintf("%s/%s", baseURL, path)
	}

	return fmt.Sprintf("http://%s/%s", baseURL, path)
}
