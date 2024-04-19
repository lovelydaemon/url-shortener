package v1

import "errors"

var (
	ErrConflict            = errors.New("data conflict")
	ErrInternalServerError = errors.New("internal server error")
	ErrNotFound            = errors.New("record not found")
	ErrNoUpdate            = errors.New("no update")
)
