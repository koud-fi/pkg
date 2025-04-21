package flex

import "encoding/json"

// OneOrMany is a generic type that can hold either a single T or a slice of T.
// This is useful when parsing data formats (like JSON or YAML) where a field might
// be represented as a single value or as a list.
type OneOrMany[T any] []T

// One creates a OneOrMany with a single value.
func One[T any](value T) OneOrMany[T] {
	return OneOrMany[T]{value}
}

// Many creates a OneOrMany from multiple values.
func Many[T any](values ...T) OneOrMany[T] {
	return OneOrMany[T](values)
}

// Value return the first value of the underlying slice, or a zero value.
func (o OneOrMany[T]) Value() T {
	if len(o) == 0 {
		var zero T
		return zero
	}
	return o[0]
}

func (o *OneOrMany[T]) UnmarshalJSON(data []byte) error {
	var single T
	if err := json.Unmarshal(data, &single); err == nil {
		*o = []T{single}
		return nil
	}
	var many []T
	if err := json.Unmarshal(data, &many); err != nil {
		return err
	}
	*o = many
	return nil
}

func (o OneOrMany[T]) MarshalJSON() ([]byte, error) {
	if len(o) == 1 {
		return json.Marshal(o[0])
	}
	return json.Marshal([]T(o))
}

func (o *OneOrMany[T]) UnmarshalYAML(unmarshal func(any) error) error {
	var single T
	if err := unmarshal(&single); err == nil {
		*o = []T{single}
		return nil
	}
	var many []T
	if err := unmarshal(&many); err != nil {
		return err
	}
	*o = many
	return nil
}

func (o OneOrMany[T]) MarshalYAML() (any, error) {
	if len(o) == 1 {
		return o[0], nil
	}
	return []T(o), nil
}
