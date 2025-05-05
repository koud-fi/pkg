package kv

import (
	"context"
	"iter"
)

type (
	Pair[K comparable, V any] interface {
		Key() K
		Value() V
	}
	Reader[K comparable, V any] interface {
		// Get returns the value for the given key, ErrNotFound if the key does not exist.
		Get(context.Context, K) (V, error)
	}
	ScanReader[K comparable, V any] interface {
		Reader[K, V]
		// Scan returns a sequence of all values in the storage.
		Scan(context.Context) (iter.Seq[Pair[K, V]], func() error)
	}
	Storage[K comparable, V any] interface {
		Reader[K, V]
		// Set sets the value for the given key, overwriting any existing value.
		Set(context.Context, K, V) error
		// Del deletes the value for the given key, does nothing if the key does not exist.
		Del(context.Context, K) error
	}
)
