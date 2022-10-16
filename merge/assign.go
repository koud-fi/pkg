package merge

import (
	"fmt"
	"reflect"
	"strings"
)

//type Getter interface{ Get(string) (any, bool) }

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
