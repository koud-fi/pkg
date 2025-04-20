package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/koud-fi/pkg/serve"
)

type (
	RPCEndpoint struct {
		endpoint
	}
	RPCOutput[T any] struct {
		Ok    bool  `json:""`
		Error error `json:",omitempty"`
		Data  T     `json:",omitempty"`
	}
)

func RPC[T1, T2 any](fn func(context.Context, T1) (T2, error)) RPCEndpoint {
	e := newEndpoint(fn)

	// TODO: validate input and output types

	return RPCEndpoint{endpoint: e}
}

func (e RPCEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serve.Handle(w, r, func() (*serve.Info, error) {
		out, err := e.call(r.Context(), func(v any) error {
			return applyHTTPInput(v, r)
		})
		if err != nil {
			return nil, fmt.Errorf("call: %w", err)
		}
		var outErr error
		if !out[1].IsNil() {
			outErr = out[1].Interface().(error)
		}
		return serve.JSON(w, r, RPCOutput[any]{
			Ok:    true,
			Data:  out[0].Interface(),
			Error: outErr,
		})
	})
}
