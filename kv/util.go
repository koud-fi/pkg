package kv

import (
	"context"
	"errors"
	"fmt"
	"iter"
)

type ValueCollector[K comparable, V any] struct {
	seq   iter.Seq[Pair[K, V]]
	errFn func() error
}

func (vc *ValueCollector[K, V]) All() ([]V, error) {
	return vc.Filter(nil)
}

func (vc *ValueCollector[K, V]) Filter(pred func(V) bool) ([]V, error) {
	values := vc.buffer()
	for p := range vc.seq {
		if pred == nil || pred(p.Value()) {
			values = append(values, p.Value())
		}
	}
	if err := vc.errFn(); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	return values, nil
}

func (vc *ValueCollector[K, V]) buffer() []V {
	// We ensure that non-nil slice is returned so it's never marshaled as null.
	return make([]V, 0) // TODO: use pooled buffers with finalizers?
}

func Values[K comparable, V any](
	seq iter.Seq[Pair[K, V]], errFn func() error,
) *ValueCollector[K, V] {
	return &ValueCollector[K, V]{seq: seq, errFn: errFn}
}

// Lookup returns the first value that satisfies the predicate.
// This is essentially a full "table scan" and should be used with caution.
func Lookup[K comparable, V any](
	s ScanReader[K, V],
	ctx context.Context,
	pred func(V) bool,
) (V, error) {
	pairs, errFn := s.Scan(ctx)
	for p := range pairs {
		if pred(p.Value()) {
			return p.Value(), nil
		}
	}
	var zero V
	if err := errFn(); err != nil {
		return zero, fmt.Errorf("scan: %w", err)
	}
	return zero, ErrNotFound
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
