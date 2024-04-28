package errors

import "errors"

var (
	ErrConflict    = errors.New("data conflict")
	ErrRecNotFound = errors.New("record not found")
)
