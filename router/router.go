package router

type Router[T any] struct {
	routeMap map[string]T
}

type Route[T any] struct {
	Path    string
	Handler T
}

func New[T any]() Router[T] {
	return Router[T]{routeMap: make(map[string]T)}
}

func (r *Router[T]) Add(route string, handler T) {
	r.routeMap[route] = handler
}

func (r Router[T]) Lookup(route string) (T, func(string) string) {

	// TODO: support route parameters

	return r.routeMap[route], func(key string) string { return "" }
}

func (r Router[T]) Routes() []Route[T] {
	rs := make([]Route[T], 0, len(r.routeMap))
	for p, h := range r.routeMap {
		rs = append(rs, Route[T]{
			Path:    p,
			Handler: h,
		})
	}
	return rs
}
