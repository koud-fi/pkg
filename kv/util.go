package kv

import (
	"context"
	"errors"
	"fmt"
	"iter"
)

func Values[K comparable, V any](
	seq iter.Seq[Pair[K, V]], errFn func() error,
) ([]V, error) {
	// We ensure that non-nil slice is returned so it's never marshaled as null.
	values := make([]V, 0) // TODO: use pooled buffers with finalizers?
	for p := range seq {
		values = append(values, p.Value())
	}
	if err := errFn(); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	return values, nil
}

func Update[K comparable, V any](
	s Storage[K, V],
	ctx context.Context,
	key K,
	update func(V) (V, error),
) (V, error) {
	v, err := s.Get(ctx, key)
	if err != nil {
		return v, fmt.Errorf("get current value: %w", err)
	}
	if v, err = update(v); err != nil {
		return v, fmt.Errorf("update: %w", err)
	}
	return v, s.Set(ctx, key, v)
}

func Upsert[K comparable, V any](
	s Storage[K, V],
	ctx context.Context,
	key K,
	create func() (K, V, error),
	update func(V) (V, error),
) (v V, err error) {
	var zeroKey K
	if key != zeroKey {
		v, err = s.Get(ctx, key)
	}
	switch {
	case err != nil:
	}
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return v, fmt.Errorf("get current value: %w", err)
		}
	}
	if key == zeroKey {
		if key, v, err = create(); err != nil {
			return v, fmt.Errorf("create: %w", err)
		}
		if key == zeroKey {
			return v, fmt.Errorf("create: key is zero")
		}
	} else {
		if v, err = update(v); err != nil {
			return v, fmt.Errorf("update: %w", err)
		}
	}
	return v, s.Set(ctx, key, v)
}
