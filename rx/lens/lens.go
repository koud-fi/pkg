package lens

import (
	"sync"

	"github.com/koud-fi/pkg/rx"
)

type valueLens[T any] struct{ v T }

func (vl valueLens[T]) Get() (T, error) { return vl.v, nil }
func (vl *valueLens[T]) Set(v T) error  { vl.v = v; return nil }

func Value[T any](v T) rx.Lens[T] { return &valueLens[T]{v} }

type atomicLens[T any] struct {
	l  rx.Lens[T]
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

func Atomic[T any](l rx.Lens[T]) rx.Lens[T] { return &atomicLens[T]{l: l} }

type onceLens[T any] struct {
	fn   func() (T, error)
	init bool
	v    T
	err  error
}

func Once[T any](fn func() (T, error)) rx.Lens[T] { return &onceLens[T]{fn: fn} }

func (ol *onceLens[T]) Get() (T, error) {
	if !ol.init {
		v, err := ol.fn()
		ol.init, ol.v, ol.err = true, v, err
	}
	return ol.v, ol.err
}

func (ol *onceLens[T]) Set(v T) error {
	ol.init, ol.v, ol.err = true, v, nil
	return nil
}
