package realtimeforum

import "errors"

var (
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
	ErrInternal     = errors.New("internal error")
	ErrMethodNotAllowed = errors.New("method not allowed")
	ErrInvalidPassForm = errors.New("password must contain at least one letter, one number and one symbol")
)
