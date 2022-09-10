package datastore

import (
	"encoding/json"

	"github.com/koud-fi/pkg/blob"
)

type Codec[T any] interface {
	Marshal(T) blob.Blob
	Unmarshal(blob.Blob) (T, error)
}

func JSON[T any]() Codec[T] { return jsonCodec[T]{} }

type jsonCodec[T any] struct{}

func (jsonCodec[T]) Marshal(v T) blob.Blob {
	return blob.Marshal(json.Marshal, v)
}

func (jsonCodec[T]) Unmarshal(b blob.Blob) (v T, _ error) {
	return v, blob.Unmarshal(json.Unmarshal, b, &v)
}
