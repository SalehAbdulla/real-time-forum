package realtimeforum

import "errors"

var (
	ErrBadRequest         = errors.New("bad request")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrNotFound           = errors.New("not found")
	ErrInternal           = errors.New("internal error")
	ErrGender             = errors.New("invalid gender")
	ErrMethodNotAllowed   = errors.New("method not allowed")
	ErrInvalidPassForm    = errors.New("password must contain at least one letter, one number and one symbol")
	ErrInvalidCredentials = errors.New("invalid email/nickname or password")
	ErrPasswordsDontMatch = errors.New("passwords do not match")
	ErrInvalidAge         = errors.New("invalid age")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrEmailExists        = errors.New("email already exists")
	ErrNickName           = errors.New("nickname already taken")
)
