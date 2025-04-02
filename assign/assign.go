package assign

import (
	"errors"
	"reflect"
)

func Value[T any](out *T, in any) error {
	return ValueWithConverter(out, in, DefaultConverter)
}

func ValueWithConverter[T any](out *T, in any, conv *Converter) error {
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("out must be a non-nil pointer")
	}
	target := rv.Elem().Type()
	convVal, err := DefaultConverter.Convert(in, target)
	if err != nil {
		return err
	}
	rv.Elem().Set(convVal)
	return nil
}
