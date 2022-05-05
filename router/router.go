package router

type Router struct {
	routes map[string]any
}

func New() *Router {
	return &Router{routes: make(map[string]any)}
}

func (r *Router) Add(route string, handler any) {
	r.routes[route] = handler
}

func (r Router) Lookup(route string) (any, func(string) string) {

	// TODO: support route parameters

	return r.routes[route], func(key string) string { return "" }
}
