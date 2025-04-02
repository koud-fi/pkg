package assign

import "reflect"

func Implements[T any](t reflect.Type) bool {
	return t.Implements(reflect.TypeOf((*T)(nil)).Elem())
}
