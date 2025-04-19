package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/koud-fi/pkg/assign"
)

type endpoint struct {
	fn          reflect.Value
	inType      reflect.Type
	inTypeIsPtr bool
	outType     reflect.Type
}

func newEndpoint(fn any) endpoint {
	var (
		fnVal  = reflect.ValueOf(fn)
		fnType = fnVal.Type()
		e      = endpoint{
			fn:      fnVal,
			inType:  fnType.In(1),
			outType: fnType.Out(0),
		}
	)
	if e.inType.Kind() == reflect.Ptr {
		e.inType = e.inType.Elem()
		e.inTypeIsPtr = true
	}
	return e
}

func (e endpoint) call(
	ctx context.Context, applyArgs func(any) error,
) ([]reflect.Value, error) {
	args := [2]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.New(e.inType),
	}
	if err := applyArgs(args[1].Interface()); err != nil {
		return nil, fmt.Errorf("can't apply input: %w", err)
	}
	if !e.inTypeIsPtr {
		args[1] = args[1].Elem()
	}
	return e.fn.Call(args[:]), nil
}

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
