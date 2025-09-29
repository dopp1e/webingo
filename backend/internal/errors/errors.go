package errors

import "errors"

var (
	ErrDataVersionMismatch = errors.New("data version mismatch, please refresh the page or restart the application")
	ErrNotFound            = errors.New("resource not found")
	ErrAlreadyExists       = errors.New("resource already exists")
	ErrUsernameExists      = errors.New("user with this username already exists")
	ErrEmailExists         = errors.New("user with this email already exists")
	ErrTokenRequired       = errors.New("token is required")
	ErrTokenExpired        = errors.New("token has expired")
)
