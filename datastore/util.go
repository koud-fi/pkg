package datastore

import (
	"context"

	"github.com/koud-fi/pkg/rx"
)

func Sync[T any](ctx context.Context, dst, src *Sorted[T], after string) error {
	return rx.ForEach(src.Iter(ctx, after), func(p rx.Pair[string, T]) error {
		return dst.Set(ctx, p.Key, p.Value)
	})
}
