package router

type Router[T any] struct {
	// TODO
}

func (r Router[T]) Register(method, path string, v T) {

	// ???

	panic("TODO")
}

func (r Router[T]) Lookup(method, path string) (T, func(string) string, bool) {

	// ???

	panic("TODO")
}
