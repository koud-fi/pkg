package router

import "github.com/koud-fi/pkg/rx"

type Router[T any] struct {
	tree map[string]node[T]
}

func New[T any]() *Router[T] {
	return &Router[T]{tree: make(map[string]node[T])}
}

func (r *Router[T]) Register(method, path string, v T) {

	// ???

	panic("TODO")
}

func (r *Router[T]) Lookup(method, path string) (T, Params, bool) {

	// ???

	panic("TODO")
}

type Params struct {
	pairs []rx.Pair[string, string]
}

func (p Params) Get(key string) string {
	for _, pair := range p.pairs {
		if pair.Key == key {
			return pair.Value
		}
	}
	panic("router: undefined path parameter: " + key)
}

type node[T any] struct {
	node    map[string]node[T]
	isParam bool
	value   T
}
