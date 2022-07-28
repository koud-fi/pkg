package httpapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/serve"
)

type Endpoint struct {
	fn          reflect.Value
	inType      reflect.Type
	inTypeIsPtr bool
	outType     reflect.Type
	opts        []serve.Option
}

func New[T1, T2 any](fn func(context.Context, T1) (T2, error), opt ...serve.Option) Endpoint {
	var (
		fnVal  = reflect.ValueOf(fn)
		fnType = fnVal.Type()
		e      = Endpoint{
			fn:      fnVal,
			inType:  fnType.In(1),
			outType: fnType.Out(0),
			opts:    opt,
		}
	)
	if e.inType.Kind() == reflect.Ptr {
		e.inType = e.inType.Elem()
		e.inTypeIsPtr = true
	}

	// TODO: validate input and output types

	return e
}

func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serve.Handle(w, r, func() (*serve.Info, error) {
		var args [2]reflect.Value
		args[0] = reflect.ValueOf(r.Context())

		// TODO: args[1]

		out := e.fn.Call(args[:])
		if !out[1].IsNil() {
			return nil, out[1].Interface().(error)
		}
		switch v := out[0].Interface().(type) {
		case nil:
			return serve.Blob(w, r, blob.Empty(), e.opts...)
		case io.Reader:
			return serve.Reader(w, r, v, e.opts...)
		case blob.Blob:
			return serve.Blob(w, r, v, e.opts...)
		case []byte:
			return serve.Blob(w, r, blob.FromBytes(v), e.opts...)
		case string:
			return serve.Blob(w, r, blob.FromString(v), e.opts...)
		case fmt.Stringer:
			return serve.Blob(w, r, blob.FromString(v.String()), e.opts...)
		default:
			return serve.JSON(w, r, v, e.opts...)
		}
	})
}
