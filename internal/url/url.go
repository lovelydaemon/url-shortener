package url

import (
	"fmt"
	"strings"
)

func CreateValidURL(baseURL string, path string) string{
  if strings.HasPrefix(baseURL, "http") {
    return fmt.Sprintf("%s/%s",baseURL, path)
  }

  return fmt.Sprintf("http://%s/%s",baseURL, path)
}
