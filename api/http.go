package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/serve"
)

type HTTPOutput[T any] struct {
	Ok    bool  `json:"ok"`
	Error error `json:"error,omitempty"`
	Data  T     `json:"data,omitempty"`
}

func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serve.Handle(w, r, func() (*serve.Info, error) {
		args, err := httpRequestArgs(r)
		if err != nil {
			return nil, fmt.Errorf("http request args: %w", err)
		}
		out, err := e.Call(r.Context(), args)
		return serveOutput(w, r, out, err)
	})
}

func httpRequestArgs(r *http.Request) (Arguments, error) {
	var bodyArgs Arguments
	switch r.Header.Get("Content-Type") {
	case "application/json", "application/json; charset=UTF-8": // TODO: Make this more robust
		args := make(ArgumentMap)
		if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
			return nil, fmt.Errorf("decode json: %w", err)
		}
		bodyArgs = args
	}
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("parse form: %w", err)
	}
	args := CombinedArguments{
		URLValueArguments(r.Form),
	}
	if bodyArgs != nil {
		args = append(args, bodyArgs)
	}
	return args, nil
}

func serveOutput(
	w http.ResponseWriter, r *http.Request, output any, err error,
) (*serve.Info, error) {
	switch v := output.(type) {
	case io.ReadCloser:
		defer v.Close()
		return serve.Reader(w, r, v)
	case io.Reader:
		return serve.Reader(w, r, v)
	case blob.Blob:
		return serve.Blob(w, r, v)
	case []byte:
		return serve.Blob(w, r, blob.FromBytes(v))
	case string:
		return serve.Blob(w, r, blob.FromString(v))
	case fmt.Stringer:
		return serve.Blob(w, r, blob.FromString(v.String()))
	default:
		return serve.JSON(w, r, HTTPOutput[any]{
			Ok:    err == nil,
			Data:  v,
			Error: err,
		})
	}
}
