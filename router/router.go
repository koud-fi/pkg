package router

type Router struct {
	routeMap map[string]any
}

type Route struct {
	Path    string
	Handler any
}

func New() *Router {
	return &Router{routeMap: make(map[string]any)}
}

func (r *Router) Add(route string, handler any) {
	r.routeMap[route] = handler
}

func (r Router) Lookup(route string) (any, func(string) string) {

	// TODO: support route parameters

	return r.routeMap[route], func(key string) string { return "" }
}

func (r Router) Routes() []Route {
	rs := make([]Route, 0, len(r.routeMap))
	for p, h := range r.routeMap {
		rs = append(rs, Route{
			Path:    p,
			Handler: h,
		})
	}
	return rs
}
