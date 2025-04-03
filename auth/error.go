package auth

import "errors"

var (
	ErrBadCredentials  = errors.New("bad credentials")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrUnsupportedType = errors.New("unsupported type")
)

// TODO: proper error type
