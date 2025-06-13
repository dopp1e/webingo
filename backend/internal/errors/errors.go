package errors

import "errors"

var (
	ErrDataVersionMismatch = errors.New("data version mismatch, please refresh the page or restart the application")
	ErrNotFound            = errors.New("resource not found")
)
