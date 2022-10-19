package schema

import (
	"reflect"
)

type Schema struct {
	Type
	Definitions map[string]Type `json:"definitions,omitempty"`
}

func Resolve[T any](opts ...Option) Schema {
	c := config{
		customTypes: map[typeKey]func(reflect.Type) Type{
			{"time", "Time"}: func(_ reflect.Type) Type {
				return Type{Type: String, Format: DateTime}
			},
		},
	}
	for _, opt := range opts {
		opt(&c)
	}
	var v T
	return Schema{Type: resolveType(c, reflect.TypeOf(v))}
}
