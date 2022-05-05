package proc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Params interface {
	Apply(interface{}) error
}

type ParamFunc func(v interface{}) error

func (fn ParamFunc) Apply(v interface{}) error { return fn(v) }

type ParamMap map[string][]string

func (m ParamMap) Apply(v interface{}) error {
	return applyParams(reflect.ValueOf(v), func(key string) []string {
		return m[key]
	})
}

func applyParams(v reflect.Value, lookup func(key string) []string) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// TODO: non-struct value handling

	var (
		vType = v.Type()
		vLen  = v.NumField()
	)
	for i := 0; i < vLen; i++ {
		var (
			f     = vType.Field(i)
			fVal  = v.Field(i)
			fName = resolveFieldName(f)
		)
		if f.Anonymous {
			applyParams(fVal, lookup)
			continue
		}
		vs := lookup(fName)
		if len(vs) == 0 {
			continue
		}
		if fVal.CanAddr() && fVal.CanInterface() {
			v := fVal.Addr().Interface()
			if u, ok := v.(json.Unmarshaler); ok {
				u.UnmarshalJSON([]byte(vs[0])) // TODO: handle error
				continue
			}
		}
		fKind := f.Type.Kind()
		switch fKind {
		case reflect.String:
			fVal.SetString(vs[0])

		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			i, _ := strconv.ParseInt(vs[0], 10, f.Type.Bits()) // TODO: handle error
			fVal.SetInt(i)

		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			u, _ := strconv.ParseUint(vs[0], 10, f.Type.Bits()) // TODO: handle error
			fVal.SetUint(u)

		case reflect.Float64, reflect.Float32:
			f, _ := strconv.ParseFloat(vs[0], f.Type.Bits()) // TODO: handle error
			fVal.SetFloat(f)

		case reflect.Bool:
			b, _ := strconv.ParseBool(vs[0]) // TODO: handle error
			fVal.SetBool(b)

		// TODO: slice handling

		// TODO: struct handling

		default:
			return fmt.Errorf("unsupported kind %v for arg: %s", fKind, fName)
		}
	}
	return nil
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
