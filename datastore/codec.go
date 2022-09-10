package datastore

import (
	"encoding/json"

	"github.com/koud-fi/pkg/blob"
)

type Codec[T any] interface {
	Marshal(T) blob.Blob
	Unmarshal(blob.Blob) (T, error)
}

var _ Codec[any] = (*Funcs[any])(nil)

type Funcs[T any] struct {
	M func(T) blob.Blob
	U func(blob.Blob) (T, error)
}

func (c Funcs[T]) Marshal(v T) blob.Blob            { return c.M(v) }
func (c Funcs[T]) Unmarshal(b blob.Blob) (T, error) { return c.U(b) }

func JSON[T any]() Codec[T] {
	return Funcs[T]{
		M: func(v T) blob.Blob { return blob.Marshal(json.Marshal, v) },
		U: func(b blob.Blob) (v T, _ error) { return v, blob.Unmarshal(json.Unmarshal, b, &v) },
	}
}

func Strings() Codec[string] {
	return Funcs[string]{M: blob.FromString, U: blob.String}
}

func Bytes() Codec[[]byte] {
	return Funcs[[]byte]{M: blob.FromBytes, U: blob.Bytes}
}
