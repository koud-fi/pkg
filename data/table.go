package data

import (
	"context"

	"github.com/koud-fi/pkg/rx"
)

type Table[T any] interface {
	Get(ctx context.Context, keys rx.Iter[T]) rx.Iter[rx.Pair[T, rx.Maybe[T]]]
	Put(ctx context.Context, values rx.Iter[T]) rx.Iter[T]
	Delete(ctx context.Context, keys rx.Iter[T]) error
}

// TODO: SortedTable interface
