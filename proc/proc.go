package proc

import (
	"context"
	"fmt"
	"reflect"
)

var (
	ctxType = reflect.TypeOf((*context.Context)(nil)).Elem()
	errType = reflect.TypeOf((*error)(nil)).Elem()
)

type ctx = context.Context

type Func[T1, T2 any] interface {
	func(T1) T2 | func(ctx, T1) T2 | func(T1) (T2, error) | func(ctx, T1) (T2, error)
}

type InFunc[T any] interface {
	func(T) | func(ctx, T) | func(T) error | func(ctx, T) error
}

type OutFunc[T any] interface {
	func() T | func(ctx) T | func() (T, error) | func(ctx) (T, error)
}

type NullFunc interface {
	func() | func() error | func(ctx) | func(ctx) error
}

type Proc struct {
	fn          reflect.Value
	numIn       int
	useCtx      bool
	inType      reflect.Type
	inTypeIsPtr bool
	outType     reflect.Type
	useErr      bool
}

func New[T1, T2 any, F Func[T1, T2]](fn F) Proc { return newProc(fn) }
func NewIn[T any, F InFunc[T]](fn F) Proc       { return newProc(fn) }
func NewOut[T any, F OutFunc[T]](fn F) Proc     { return newProc(fn) }
func NewNull[F NullFunc](fn F) Proc             { return newProc(fn) }

func newProc(fn any) (pr Proc) {
	pr.fn = reflect.ValueOf(fn)
	fnTyp := pr.fn.Type()
	pr.numIn = fnTyp.NumIn()
	switch pr.numIn {
	case 1:
		if fnTyp.In(0).Implements(ctxType) {
			pr.useCtx = true
		} else {
			pr.inType = fnTyp.In(0)
		}
	case 2:
		pr.useCtx = true
		pr.inType = fnTyp.In(1)
	}
	if pr.inType != nil {
		if pr.inType.Kind() == reflect.Ptr {
			pr.inType = pr.inType.Elem()
			pr.inTypeIsPtr = true
		}
	}
	switch fnTyp.NumOut() {
	case 1:
		if fnTyp.Out(0).Implements(errType) {
			pr.useErr = true
		} else {
			pr.outType = fnTyp.Out(0)
		}
	case 2:
		pr.outType = fnTyp.Out(0)
		pr.useErr = true
	}
	return pr
}

func (pr Proc) Invoke(ctx context.Context, p Params) (any, error) {
	var callArgs []reflect.Value
	if pr.numIn > 0 {
		callArgs = make([]reflect.Value, 0, pr.numIn)
		if pr.useCtx {
			callArgs = append(callArgs, reflect.ValueOf(ctx))
		}
		if pr.inType != nil {
			inVal := reflect.New(pr.inType)
			if err := p.Apply(inVal.Interface()); err != nil {
				return nil, fmt.Errorf("can't apply params: %w", err)
			}
			if !pr.inTypeIsPtr {
				inVal = inVal.Elem()
			}
			callArgs = append(callArgs, inVal)
		}
	}
	out := pr.fn.Call(callArgs)
	switch len(out) {
	case 0:
		return nil, nil
	case 1:
		if pr.useErr {
			return nil, convertErr(out[0])
		}
		return out[0].Interface(), nil
	default:
		return out[0].Interface(), convertErr(out[1])
	}
}

func convertErr(errVal reflect.Value) error {
	if errVal.IsNil() {
		return nil
	}
	return errVal.Interface().(error)
}
