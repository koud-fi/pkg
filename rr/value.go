package rr

import (
	"reflect"
	"strconv"
)

// TODO: []byte
// TODO: bool
// TODO: []byte set
// TODO: list (slice)
// TODO: map
// TODO: number set
// TODO: null
// TODO: string
// TODO: string set

type V struct {
	s string
	//m map[string]V
}

func Int[T interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}](v T) V {
	return V{s: strconv.FormatInt(reflect.ValueOf(v).Int(), 10)}
}

func Float[T ~float32 | ~float64](v T) V {
	return V{s: strconv.FormatFloat(reflect.ValueOf(v).Float(), 'f', -1, reflect.TypeOf(v).Bits())}
}
