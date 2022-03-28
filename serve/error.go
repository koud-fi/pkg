package serve

import (
	"errors"
	"io/fs"
	"net/http"
)

type StatusCoder interface {
	StatusCode() int
}

func Error(w http.ResponseWriter, err error) {
	switch {
	case err == nil:
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, err.Error(), ErrorStatusCode(err))
	}
}

func ErrorStatusCode(err error) int {
	if scer, ok := err.(StatusCoder); ok {
		return scer.StatusCode()
	}
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return http.StatusNotFound
	case errors.Is(err, fs.ErrPermission):
		return http.StatusForbidden
	case errors.Is(err, fs.ErrInvalid):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
