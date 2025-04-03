package oauth

import (
	"net/url"
)

type CallbackError struct {
	Code        string `json:"error"`
	Description string `json:"error_description"`
}

func parseCallbackError(q url.Values) error {
	code := q.Get("error")
	if code == "" {
		return nil
	}
	return &CallbackError{
		Code:        code,
		Description: q.Get("error_description"),
	}
}

func (ce CallbackError) Error() string {
	return ce.Code + ": " + ce.Description
}
