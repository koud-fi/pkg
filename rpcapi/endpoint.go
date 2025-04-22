package rpcapi

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/koud-fi/pkg/assert"
	"github.com/koud-fi/pkg/assign"
	"github.com/koud-fi/pkg/blob"
)

type Endpoint struct {
	fn          reflect.Value
	useCtx      bool
	inType      reflect.Type
	inTypeIsPtr bool
	outType     reflect.Type
	returnErr   bool

	converter *assign.Converter
}

// NewEndpoint inspects the signature of fn and creates an Endpoint.
func NewEndpoint(converter *assign.Converter, fn any) (*Endpoint, error) {
	var (
		fnVal  = reflect.ValueOf(fn)
		fnType = fnVal.Type()
	)
	if fnType.Kind() != reflect.Func {
		return nil, errors.New("fn must be a function")
	}
	e := Endpoint{
		converter: converter,
		fn:        fnVal,
	}
	var inPos int

	switch fnType.NumIn() {
	case 0:
		inPos = -1
	case 1:
		inPos = 0
	case 2:
		ctxType := fnType.In(0)
		if ctxType == reflect.TypeOf((*context.Context)(nil)).Elem() {
			e.useCtx = true
			inPos = 1
		} else {
			return nil, fmt.Errorf(
				"if two parameters, first must be context.Context, got %v", ctxType)
		}
	default:
		return nil, fmt.Errorf(
			"fn must take 0, 1 (arg) or 2 (context, arg) parameters, got %d", fnType.NumIn())
	}
	if inPos >= 0 {
		inType := fnType.In(inPos)

		if inType.Kind() == reflect.Ptr {
			e.inTypeIsPtr = true
			e.inType = inType.Elem()
		} else {
			e.inType = inType
		}
	}
	switch fnType.NumOut() {
	case 1:
		e.outType = fnType.Out(0)
	case 2:
		errTy := fnType.Out(1)
		if errTy == reflect.TypeOf((*error)(nil)).Elem() {
			e.returnErr = true
			e.outType = fnType.Out(0)
		} else {
			return nil, fmt.Errorf(
				"second return value must be error, got %v", errTy)
		}
	default:
		return nil, fmt.Errorf(
			"fn must return 1 or 2 values, got %d", fnType.NumOut())
	}
	return &e, nil
}

func RPC[T1, T2 any](fn func(context.Context, T1) (T2, error)) *Endpoint {
	return assert.Must(NewEndpoint(assign.NewDefaultConverter(), fn))
}

func Blob[T any](fn func(context.Context, T) blob.Reader) *Endpoint {
	return assert.Must(NewEndpoint(assign.NewDefaultConverter(), fn))
}

func (e *Endpoint) Call(ctx context.Context, args Arguments) (any, error) {
	var argVal reflect.Value
	if e.inType != nil {
		argPtr := reflect.New(e.inType)
		if e.inTypeIsPtr {
			argVal = argPtr
		} else {
			argVal = argPtr.Elem()
		}
		if err := ApplyArguments(argPtr.Interface(), e.converter, args); err != nil {
			return nil, fmt.Errorf("can't apply input: %w", err)
		}
	}
	var inputs []reflect.Value
	if e.useCtx {
		inputs = append(inputs, reflect.ValueOf(ctx))
	}
	if argVal.IsValid() {
		inputs = append(inputs, argVal)
	}
	results := e.fn.Call(inputs)

	var err error
	if e.returnErr {
		errVal := results[1]
		if !errVal.IsNil() {
			err = errVal.Interface().(error)
		}
	}
	return results[0].Interface(), err
}
