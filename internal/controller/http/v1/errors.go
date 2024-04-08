package v1

import "errors"

var (
	errInternalServerError = errors.New("Internal server error")
	errNotFound            = errors.New("Record not found")
	errRedirectBlocked     = errors.New("HTTP redirect blocked")
)
