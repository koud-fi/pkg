package api

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/koud-fi/pkg/serve"
)

type (
	RPCEndpoint struct {
		fn          reflect.Value
		inType      reflect.Type
		inTypeIsPtr bool
		outType     reflect.Type
	}
	RPCOutput[T any] struct {
		Ok    bool  `json:"ok"`
		Error error `json:"error,omitempty"`
		Data  T     `json:"data,omitempty"`
	}
)

func RPC[T1, T2 any](fn func(context.Context, T1) (T2, error)) RPCEndpoint {
	var (
		fnVal  = reflect.ValueOf(fn)
		fnType = fnVal.Type()
		e      = RPCEndpoint{
			fn:      fnVal,
			inType:  fnType.In(1),
			outType: fnType.Out(0),
		}
	)
	if e.inType.Kind() == reflect.Ptr {
		e.inType = e.inType.Elem()
		e.inTypeIsPtr = true
	}

	// TODO: validate input and output types

	return e
}

func (e RPCEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serve.Handle(w, r, func() (*serve.Info, error) {
		args := [2]reflect.Value{
			reflect.ValueOf(r.Context()),
			reflect.New(e.inType),
		}
		if err := applyHTTPInput(args[1].Interface(), r); err != nil {
			return nil, fmt.Errorf("can't apply input: %w", err)
		}
		if !e.inTypeIsPtr {
			args[1] = args[1].Elem()
		}
		var (
			out    = e.fn.Call(args[:])
			outErr error
		)
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
