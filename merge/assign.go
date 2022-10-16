package merge

import (
	"fmt"
	"reflect"
)

//type Getter interface{ Get(string) (any, bool) }

type AssignFunc func(reflect.Value, any) error

func NewAssignFunc(typ reflect.Type) AssignFunc {
	switch typ.Kind() {
	case reflect.Ptr:
		var (
			elTyp = typ.Elem()
			elFn  = NewAssignFunc(typ)
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
		panic("TODO")

	default:
		panic(fmt.Errorf("merge: unsupported kind: %v", typ.Kind()))
	}
}
