package blobserver

import (
	"context"
	"fmt"
	"io"
	"io/fs"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/pk"
)

var registry = make(map[pk.Scheme]Fetcher)

// Register calls are not thread-safe and should only be done on startup.
func Register(s pk.Scheme, f Fetcher) {
	registry[s] = f
}

type TransformFunc func(ctx context.Context, params string, ref pk.Ref) blob.Blob

func TransformFetcher(fn TransformFunc) Fetcher {
	return FetchFunc(func(ctx context.Context, ref pk.Ref) blob.Blob {
		return blob.Func(func() (io.ReadCloser, error) {
			dataRef, err := pk.ParseRef(ref.Key())
			if err != nil {
				return nil, err
			}
			return fn(ctx, ref.Params(), dataRef).Open()
		})
	})
}

func Fetch(ctx context.Context, ref pk.Ref) blob.Blob {
	return Lookup(ref.Scheme()).Fetch(ctx, ref)
}

func Lookup(s pk.Scheme) Fetcher {
	f, ok := registry[s]
	if !ok {
		return FetchFunc(notFoundFetch)
	}
	return f
}

func notFoundFetch(_ context.Context, ref pk.Ref) blob.Blob {
	return blob.Func(func() (io.ReadCloser, error) {
		return nil, fmt.Errorf("%w: unknown scheme; %v", fs.ErrNotExist, ref.Scheme())
	})
}
