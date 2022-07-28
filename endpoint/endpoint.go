package httpapi

import (
	"context"
	"net/http"
	"reflect"
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
	return e
}

func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// TODO

}
