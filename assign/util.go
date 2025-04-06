package assign

import "reflect"

func Implements[T any](t reflect.Type) bool {
	target := reflect.TypeOf((*T)(nil)).Elem()
	return t.Implements(target) || reflect.PointerTo(t).Implements(target)
}
