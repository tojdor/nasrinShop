package myerrors

import "errors"

var (
	ErrNotFound = errors.New("Not Found 404")
	ErrBadRequest = errors.New("Bad request")
)
