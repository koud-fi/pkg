package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"

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
		args := [2]reflect.Value{
			reflect.ValueOf(r.Context()),
			reflect.New(e.inType),
		}
		if err := applyInput(args[1].Interface(), r); err != nil {
			return nil, fmt.Errorf("can't apply params: %w", err)
		}
		if !e.inTypeIsPtr {
			args[1] = args[1].Elem()
		}
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

func applyInput(v any, r *http.Request) error {
	ct := r.Header.Get("Content-Type")
	switch {
	case strings.HasPrefix("application/json", ct):
		return json.NewDecoder(r.Body).Decode(v)

	// TODO: support multipart forms

	default:
		if err := r.ParseForm(); err != nil {
			return err
		}
		return applyValues(reflect.ValueOf(v), func(key string) []string {
			return r.Form[key]
		})
	}
}

func applyValues(v reflect.Value, lookup func(key string) []string) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	var (
		vType = v.Type()
		vLen  = v.NumField()
		err   error
	)
	for i := 0; i < vLen && err == nil; i++ {
		var (
			f     = vType.Field(i)
			fVal  = v.Field(i)
			fName = resolveFieldName(f)
		)
		if f.Anonymous {
			applyValues(fVal, lookup)
			continue
		}
		vs := lookup(fName)
		if len(vs) == 0 {
			continue
		}
		if fVal.CanAddr() && fVal.CanInterface() {
			vi := fVal.Addr().Interface()
			if u, ok := vi.(json.Unmarshaler); ok {
				err = u.UnmarshalJSON([]byte(vs[0]))
				continue
			}
		}
		fKind := f.Type.Kind()
		switch fKind {
		case reflect.String:
			fVal.SetString(vs[0])

		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			var i int64
			if i, err = strconv.ParseInt(vs[0], 10, f.Type.Bits()); err == nil {
				fVal.SetInt(i)
			}
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			var u uint64
			if u, err = strconv.ParseUint(vs[0], 10, f.Type.Bits()); err == nil {
				fVal.SetUint(u)
			}
		case reflect.Float64, reflect.Float32:
			var fl float64
			if fl, err = strconv.ParseFloat(vs[0], f.Type.Bits()); err == nil {
				fVal.SetFloat(fl)
			}
		case reflect.Bool:
			var b bool
			if b, err = strconv.ParseBool(vs[0]); err == nil {
				fVal.SetBool(b)
			}

		// TODO: slice handling

		// TODO: struct handling

		default:
			return fmt.Errorf("unsupported kind %v for value: %s", fKind, fName)
		}
	}
	return err
}

func resolveFieldName(f reflect.StructField) string {
	name := f.Name
	if tag, ok := f.Tag.Lookup("json"); ok {
		parts := strings.SplitN(tag, ",", 2)
		if parts[0] != "" {
			name = parts[0]
		}
	}
	return name
}
