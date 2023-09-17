package datastore

import (
	"context"

	"github.com/koud-fi/pkg/rx"
)

type KV[T any] interface {
	Get(ctx context.Context, key string) (T, error)
	Put(ctx context.Context, key string, value T) error
	Delete(ctx context.Context, keys ...string) error
}

type SortedKV[T any] interface {
	KV[T]
	Iter(ctx context.Context, state rx.Lens[string]) rx.Iter[rx.Pair[string, T]]
}
