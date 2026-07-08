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
	ErrInvalidCredentials = errors.New("invalid email/username or password")
	ErrPasswordsDontMatch = errors.New("passwords do not match")
	ErrInvalidAge         = errors.New("invalid age")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrEmailExists        = errors.New("email already exists")
	ErrNickName           = errors.New("username already taken")
	ErrNickNameLength = errors.New("invalid username length (min 2, max 33)")
	ErrPasswordLength = errors.New("invalid password length (min 12, max 64)")
	ErrTitleEmptyOrMoreThanHundard = errors.New("invalid title length, (min 3, max 30) characters")
	ErrContentEmptyOrMoreThanHundard = errors.New("invalid content length, (min 10, max 100) characters")
	ErrNoCategorySelected = errors.New("post must have a valid category")
	ErrMissingPostId = errors.New("missing post id")
)
