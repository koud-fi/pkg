package rpcapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/serve"
)

type HTTPOutput[T any] struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
	Data  T      `json:"data,omitzero"`
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serve.Handle(w, r, func() (*serve.Info, error) {
		args, err := httpRequestArgs(r)
		if err != nil {
			return nil, fmt.Errorf("http request args: %w", err)
		}
		out, err := e.Call(r.Context(), args)
		return serveOutput(w, r, out, err)
	})
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serve.Handle(w, r, func() (*serve.Info, error) {
		args, err := httpRequestArgs(r)
		if err != nil {
			return nil, fmt.Errorf("http request args: %w", err)
		}
		pathParts := strings.Split(r.URL.Path, "/")
		name := pathParts[len(pathParts)-1]

		out, err := m.Call(r.Context(), name, args)
		return serveOutput(w, r, out, err)
	})
}

func httpRequestArgs(r *http.Request) (Arguments, error) {
	var (
		bodyArgs    Arguments
		contentType = r.Header.Get("Content-Type")
	)
	switch {
	case contentType == "":
		break // no body

	case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"):
		break // form is parsed later (to also handle possible query params)

	case strings.HasPrefix(contentType, "application/json"):
		args := make(ArgumentMap)
		if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
			return nil, fmt.Errorf("decode json: %w", err)
		}
		bodyArgs = args
	default:
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
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
		var errStr string
		if err != nil {
			errStr = err.Error()
		}
		return serve.JSON(w, r, HTTPOutput[any]{
			Ok:    err == nil,
			Data:  v,
			Error: errStr,
		})
	}
}
