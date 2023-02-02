package rx

import "sync"

type Lens[T any] interface {
	Get() (T, error)
	Set(T) error
}

type valueLens[T any] struct{ v T }

func (vl valueLens[T]) Get() (T, error) { return vl.v, nil }
func (vl *valueLens[T]) Set(v T) error  { vl.v = v; return nil }

func Value[T any](v T) Lens[T] { return &valueLens[T]{v} }

type atomicLens[T any] struct {
	l  Lens[T]
	mu sync.RWMutex
}

func (al *atomicLens[T]) Get() (T, error) {
	al.mu.RLock()
	defer al.mu.RUnlock()
	return al.l.Get()
}

func (al *atomicLens[T]) Set(v T) error {
	al.mu.Lock()
	defer al.mu.Unlock()
	return al.l.Set(v)
}

func Atomic[T any](l Lens[T]) Lens[T] { return &atomicLens[T]{l: l} }
