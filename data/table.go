package data

import (
	"context"

	"github.com/koud-fi/pkg/rx"
)

type Table[T any] interface {
	Get(ctx context.Context) func(key T) (rx.Pair[T, rx.Maybe[T]], error)
	Put(ctx context.Context) func(value T) (T, error)
	Delete(ctx context.Context) func(key T) error
}

/*
type SortedTable[T any] interface {
	Table[T]
	Iter(ctx context.Context) // TODO: ???
}
*/
