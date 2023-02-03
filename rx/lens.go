package rx

import (
	"context"
	"errors"
	"sync"
)

type Lens[T any] interface {
	Get(context.Context) (T, error)
	Set(context.Context, T) error
}

type valueLens[T any] struct{ v T }

func (vl valueLens[T]) Get(context.Context) (T, error)    { return vl.v, nil }
func (vl *valueLens[T]) Set(_ context.Context, v T) error { vl.v = v; return nil }

func Value[T any](v T) Lens[T] { return &valueLens[T]{v} }

type atomicLens[T any] struct {
	l  Lens[T]
	mu sync.RWMutex
}

func (al *atomicLens[T]) Get(ctx context.Context) (T, error) {
	al.mu.RLock()
	defer al.mu.RUnlock()
	return al.l.Get(ctx)
}

func (al *atomicLens[T]) Set(ctx context.Context, v T) error {
	al.mu.Lock()
	defer al.mu.Unlock()
	return al.l.Set(ctx, v)
}

func Atomic[T any](l Lens[T]) Lens[T] { return &atomicLens[T]{l: l} }

type onceLens[T any] struct {
	fn   func(context.Context) (T, error)
	init bool
	v    T
	err  error
}

func Once[T any](fn func(context.Context) (T, error)) Lens[T] { return &onceLens[T]{fn: fn} }

func (ol *onceLens[T]) Get(ctx context.Context) (T, error) {
	if !ol.init {
		v, err := ol.fn(ctx)
		if err != nil &&
			(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
			return v, err
		}
		ol.init, ol.v, ol.err = true, v, err
	}
	return ol.v, ol.err
}

func (ol *onceLens[T]) Set(_ context.Context, v T) error {
	ol.init, ol.v, ol.err = true, v, nil
	return nil
}
