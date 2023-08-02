package schema

import (
	"reflect"
)

type Schema struct {
	Ref string `json:"$ref,omitempty"`
	Type
	Definitions map[string]Type `json:"definitions,omitempty"`
}

func Resolve[T any](opts ...Option) Schema {
	var v T
	return ResolveFromValue(v, opts...)
}

func ResolveFromValue(v any, opts ...Option) Schema {
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
	return Schema{Type: Type{}.resolve(c, v)}
}
