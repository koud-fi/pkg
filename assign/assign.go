package assign

import (
	"reflect"

	"github.com/koud-fi/pkg/errx"
)

var defaultConverter = NewDefaultConverter()

func Value[T any](out *T, in any) error {
	return ValueWithConverter(out, in, defaultConverter)
}

func ValueWithConverter[T any](out *T, in any, conv *Converter) error {
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errx.New("out must be a non-nil pointer")
	}
	target := rv.Elem().Type()
	convVal, err := defaultConverter.Convert(in, target)
	if err != nil {
		return errx.E(err)
	}
	rv.Elem().Set(convVal)
	return nil
}
