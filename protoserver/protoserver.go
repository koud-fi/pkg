package protoserver

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/koud-fi/pkg/pk"
)

var registry = make(map[pk.Scheme]Fetcher)

type Fetcher interface {
	Fetch(context.Context, pk.Ref) (any, error)
}

type FetchFunc func(context.Context, pk.Ref) (any, error)

func (fn FetchFunc) Fetch(ctx context.Context, ref pk.Ref) (any, error) {
	return fn(ctx, ref)
}

func Register(s pk.Scheme, f Fetcher) {
	registry[s] = f
}

type TransformFunc func(ctx context.Context, params string, ref pk.Ref) (any, error)

func TransformFetcher(fn TransformFunc) Fetcher {
	return FetchFunc(func(ctx context.Context, ref pk.Ref) (any, error) {
		keyRef, err := pk.ParseRef(ref.Key())
		if err != nil {
			return nil, err
		}
		return fn(ctx, ref.Params(), keyRef)
	})
}

func Fetch(ctx context.Context, ref pk.Ref) (any, error) {
	return Lookup(ref.Scheme()).Fetch(ctx, ref)
}

func Lookup(s pk.Scheme) Fetcher {
	f, ok := registry[s]
	if !ok {
		return FetchFunc(notFoundFetch)
	}
	return f
}

func notFoundFetch(_ context.Context, ref pk.Ref) (any, error) {
	return nil, fmt.Errorf("%w: unknown scheme; %v", fs.ErrNotExist, ref.Scheme())
}
