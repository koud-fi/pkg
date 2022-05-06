package router

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/koud-fi/pkg/proc"
)

const endpointPattern = `^[\.a-zA-Z0-9]+$`

var endpointValidator = regexp.MustCompile(endpointPattern)

type Router struct {
	endpoints map[string]proc.Proc
}

func New() Router {
	return Router{endpoints: make(map[string]proc.Proc)}
}

func (r *Router) Add(endpoint string, pr proc.Proc) {
	if !endpointValidator.MatchString(endpoint) {
		panic(fmt.Sprintf("%s doesn't match endpoint pattern: %s",
			endpoint, endpointPattern))
	}
	r.endpoints[normalizeEndpoint(endpoint)] = pr
}

func (r Router) Invoke(ctx context.Context, route string, p proc.Params) (any, error) {
	pr, ok := r.endpoints[normalizeEndpoint(route)]
	if !ok {
		return nil, fmt.Errorf("not found: %s", route)
	}
	return pr.Invoke(ctx, p)
}

func (r Router) Endpoints() []string {
	es := make([]string, 0, len(r.endpoints))
	for e := range r.endpoints {
		es = append(es, e)
	}
	sort.Strings(es)
	return es
}

func normalizeEndpoint(endpoint string) string {
	return strings.ToLower(strings.Trim(endpoint, "/"))
}
