package validation

import "net/url"

func IsValidUrl(u string) error {
  _, err := url.ParseRequestURI(u)
  return err
}
