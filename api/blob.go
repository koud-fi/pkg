package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/serve"
)

type BlobEndpoint struct {
	endpoint
}

func Blob[T any](fn func(context.Context, T) blob.Reader) BlobEndpoint {
	e := newEndpoint(fn)

	// TODO: validate input and output types

	return BlobEndpoint{endpoint: e}
}

func (e BlobEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serve.Handle(w, r, func() (*serve.Info, error) {
		out, err := e.call(r.Context(), func(v any) error {
			return applyHTTPInput(v, r)
		})
		if err != nil {
			return nil, fmt.Errorf("call: %w", err)
		}
		return serve.Blob(w, r, out[0].Interface().(blob.Reader))
	})
}
