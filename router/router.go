package router

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

func (r *Router[T]) Lookup(method, path string) (T, func(string) string, bool) {

	// ???

	panic("TODO")
}

type node[T any] struct {
	node    map[string]node[T]
	isParam bool
	value   T
}
