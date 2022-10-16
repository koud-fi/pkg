package merge

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type AssignFunc func(reflect.Value, any) error

func NewAssignFunc(typ reflect.Type) AssignFunc {
	switch typ.Kind() {
	case reflect.Ptr:
		var (
			elTyp = typ.Elem()
			elFn  = NewAssignFunc(elTyp)
		)
		return func(dst reflect.Value, v any) error {
			if dst.IsZero() {
				dst.Set(reflect.New(elTyp))
			}
			return elFn(dst.Elem(), v)
		}
	case reflect.String:
		return func(dst reflect.Value, v any) error {
			switch v := v.(type) {
			case string:
				dst.SetString(v)
			case []byte:
				dst.SetString(string(v))
			default:
				dst.SetString(fmt.Sprint(v))
			}
			return nil
		}
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return func(dst reflect.Value, v any) (err error) {
			var n int64
			switch v := v.(type) {
			case string:
				n, err = strconv.ParseInt(v, 10, typ.Bits())
			case []byte:
				n, err = strconv.ParseInt(string(v), 10, typ.Bits())
			case int8:
				n = int64(v)
			case int16:
				n = int64(v)
			case int32:
				n = int64(v)
			case int64:
				n = v
			case int:
				n = int64(v)
			case float32:
				n = int64(v)
			case float64:
				n = int64(v)
			default:
				return fmt.Errorf("merge: cannot assign %T to int", v)
			}
			dst.SetInt(n)
			return
		}
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return func(dst reflect.Value, v any) (err error) {
			var n uint64
			switch v := v.(type) {
			case string:
				n, err = strconv.ParseUint(v, 10, typ.Bits())
			case []byte:
				n, err = strconv.ParseUint(string(v), 10, typ.Bits())
			case uint8:
				n = uint64(v)
			case uint16:
				n = uint64(v)
			case uint32:
				n = uint64(v)
			case uint64:
				n = v
			case uint:
				n = uint64(v)
			default:
				return fmt.Errorf("merge: cannot assign %T to uint", v)
			}
			dst.SetUint(n)
			return
		}
	case reflect.Float64, reflect.Float32:
		return func(dst reflect.Value, v any) (err error) {
			var n float64
			switch v := v.(type) {
			case string:
				n, err = strconv.ParseFloat(v, typ.Bits())
			case []byte:
				n, err = strconv.ParseFloat(string(v), typ.Bits())
			case float32:
				n = float64(v)
			case float64:
				n = v
			default:
				return fmt.Errorf("merge: cannot assign %T to float", v)
			}
			dst.SetFloat(n)
			return
		}
	case reflect.Bool:
		return func(dst reflect.Value, v any) (err error) {
			var b bool
			switch v := v.(type) {
			case bool:
				b = v
			case
				int8, int16, int32, int64, int,
				uint8, uint16, uint32, uint64, uint,
				float32, float64:

				b = v != 0
			case string:
				b, err = strconv.ParseBool(v)
			case []byte:
				b, err = strconv.ParseBool(string(v))
			default:
				return fmt.Errorf("merge: cannot assign %T to bool", v)
			}
			dst.SetBool(b)
			return
		}
	case reflect.Slice:
		var (
			elTyp = typ.Elem()
			elFn  = NewAssignFunc(elTyp)
		)
		return func(dst reflect.Value, v any) (err error) {
			vVal := reflect.ValueOf(v)
			switch vVal.Kind() {
			case reflect.Slice, reflect.Array:
				l := vVal.Len()
				if dst.Cap() < l {
					dst.Set(reflect.MakeSlice(dst.Type(), l, l))
				}
				dst.SetLen(l)
				for i := 0; i < l; i++ {
					if err := elFn(dst.Index(i), vVal.Index(i).Interface()); err != nil {
						return err
					}
				}
			default:
				return fmt.Errorf("merge: cannot assign %T to slice", v)
			}
			return nil
		}
	case reflect.Struct:
		fieldFns := make(map[string]AssignFunc, typ.NumField())
		for i := 0; i < typ.NumField(); i++ {
			var (
				idx = i
				f   = typ.Field(idx)
			)
			if !f.IsExported() {
				continue
			}
			key := strings.ToLower(f.Name)
			if _, ok := fieldFns[key]; ok {
				panic("merge: duplicate field key: " + key)
			}
			fieldFn := NewAssignFunc(f.Type)
			fieldFns[key] = func(dst reflect.Value, v any) error {
				return fieldFn(dst.Field(idx), v)
			}
		}
		return func(dst reflect.Value, v any) error {
			var (
				vVal = reflect.ValueOf(v)
				vTyp = vVal.Type()
			)
			switch vVal.Kind() {
			case reflect.Struct:
				for i := 0; i < vVal.NumField(); i++ {
					f := vTyp.Field(i)
					if !f.IsExported() {
						continue
					}
					fn, ok := fieldFns[strings.ToLower(f.Name)]
					if !ok {
						continue
					}
					if err := fn(dst, vVal.Field(i).Interface()); err != nil {
						return err
					}
				}
				return nil

			case reflect.Map:
				// TODO
				fallthrough

			default:
				return fmt.Errorf("merge: cannot assign %T to struct", v)
			}
		}
	default:
		panic(fmt.Errorf("merge: unsupported kind: %v", typ.Kind()))
	}
}
