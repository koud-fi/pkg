package blob

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/koud-fi/pkg/rx"
)

type Storage interface {
	Getter
	Setter
	Deleter
}

type SortedStorage interface {
	Storage
	Iterator
}

type Getter interface {
	Get(ctx context.Context, ref Ref) Blob
}

type Setter interface {
	Set(ctx context.Context, ref Ref, r io.Reader) error
}

type Deleter interface {
	Delete(ctx context.Context, refs ...Ref) error
}

type RefBlob struct {
	Ref Ref
	Blob
}

type Iterator interface {
	Iter(ctx context.Context, d Domain, after Ref) rx.Iter[RefBlob]
}

type GetterFunc func(ctx context.Context, ref Ref) Blob

func (f GetterFunc) Get(ctx context.Context, ref Ref) Blob { return f(ctx, ref) }

func FSGetter(fsys fs.FS) Getter {
	return GetterFunc(func(_ context.Context, ref Ref) Blob {
		return FromFS(fsys, path.Clean(strings.Join(ref, "/")))
	})
}

func Mapper(g Getter, fn func(io.ReadCloser) (io.ReadCloser, error)) Getter {
	return GetterFunc(func(ctx context.Context, ref Ref) Blob {
		return Func(func() (io.ReadCloser, error) {
			rc, err := g.Get(ctx, ref).Open()
			if err != nil {
				return nil, err
			}
			return fn(rc)
		})
	})
}

type Mux map[Domain]Getter

func (m Mux) Get(ctx context.Context, ref Ref) Blob {
	return Func(func() (io.ReadCloser, error) {
		g, err := m.Lookup(ref)
		if err != nil {
			return nil, err
		}
		return g.Get(ctx, ref).Open()
	})
}

func (m Mux) Lookup(ref Ref) (Getter, error) {
	g, ok := m[ref.Domain()]
	if !ok {
		if g, ok = m[Default]; !ok {
			return nil, os.ErrNotExist
		}
	}
	return g, nil
}
