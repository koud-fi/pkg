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
	return Schema{Type: resolveType(makeConfig(opts), reflect.TypeOf(v))}
}

func ResolveFromValue(v any, opts ...Option) Schema {
	return Schema{Type: resolveTypeFromValue(makeConfig(opts), v)}
}

func makeConfig(opts []Option) config {
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
	return c
}
