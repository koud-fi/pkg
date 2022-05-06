package proc

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

type Router struct {
	routeMap map[string]Proc
}

func NewRouter() Router {
	return Router{routeMap: make(map[string]Proc)}
}

func (r *Router) Add(route string, pr Proc) {
	r.routeMap[normalizeRoute(route)] = pr
}

func (r Router) Invoke(ctx context.Context, route string, p Params) (any, error) {
	pr, ok := r.routeMap[normalizeRoute(route)]
	if !ok {
		return nil, fmt.Errorf("not found: %s", route)
	}
	return pr.Invoke(ctx, p)
}

func (r Router) Routes() []string {
	rs := make([]string, 0, len(r.routeMap))
	for r := range r.routeMap {
		rs = append(rs, r)
	}
	sort.Strings(rs)
	return rs
}

func normalizeRoute(route string) string {
	return strings.ToLower(strings.Trim(route, "/"))
}
