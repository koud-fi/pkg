package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/koud-fi/pkg/assign"
)

func applyHTTPInput(v any, r *http.Request) error {
	var bodyArgs Arguments
	switch r.Header.Get("Content-Type") {
	case "application/json", "application/json; charset=UTF-8":
		args := make(ArgumentMap)
		if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
			return fmt.Errorf("decode json: %w", err)
		}
		bodyArgs = args
	}
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("parse form: %w", err)
	}
	args := CombinedArguments{
		URLValueArguments(r.Form),
	}
	if bodyArgs != nil {
		args = append(args, bodyArgs)
	}
	return ApplyArguments(v, assign.NewDefaultConverter(), args)
}
