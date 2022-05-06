package proc

import (
	"context"
	"fmt"

	"github.com/koud-fi/pkg/router"
)

type Router struct {
	router.Router[Proc]
}

func NewRouter() *Router {
	return &Router{Router: router.New[Proc]()}
}

func (r Router) Invoke(ctx context.Context, route string, p Params) (any, error) {
	pr, _ := r.Lookup(route)
	if !pr.IsValid() {
		return nil, fmt.Errorf("not found: %s", route)
	}

	// TODO: handle route params

	return pr.Invoke(ctx, p)
}
