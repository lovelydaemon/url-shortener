package v1

import "errors"

var ErrConflict = errors.New("data conflict")
var ErrInternalServerError = errors.New("internal server error")
var ErrNotFound = errors.New("record not found")
