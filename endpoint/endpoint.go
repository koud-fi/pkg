package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/koud-fi/pkg/assign"
	"github.com/koud-fi/pkg/serve"
)

type Endpoint struct {
	fn          reflect.Value
	inType      reflect.Type
	inTypeIsPtr bool
	outType     reflect.Type
}

func New[T1, T2 any](fn func(context.Context, T1) (T2, error)) Endpoint {
	var (
		fnVal  = reflect.ValueOf(fn)
		fnType = fnVal.Type()
		e      = Endpoint{
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

func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serve.Handle(w, r, func() (*serve.Info, error) {
		args := [2]reflect.Value{
			reflect.ValueOf(r.Context()),
			reflect.New(e.inType),
		}
		if err := applyInput(args[1].Interface(), r); err != nil {
			return nil, fmt.Errorf("can't apply input: %w", err)
		}
		if !e.inTypeIsPtr {
			args[1] = args[1].Elem()
		}
		var (
			out    = e.fn.Call(args[:])
			outVal = out[0].Interface()
			outErr error
		)
		if !out[1].IsNil() {
			outErr = out[1].Interface().(error)
		}
		return serveOutput(w, r, outVal, outErr)
	})
}

func applyInput(v any, r *http.Request) error {
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
