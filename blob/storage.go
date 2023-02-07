package blob

import (
	"context"
	"io"
	"io/fs"
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
	Ref string
	Blob
}

type Iterator interface {
	Iter(ctx context.Context, after string) rx.Iter[RefBlob]
}

type GetterFunc func(ctx context.Context, ref Ref) Blob

func (f GetterFunc) Get(ctx context.Context, ref Ref) Blob { return f(ctx, ref) }

func FSGetter(fsys fs.FS) Getter {
	return GetterFunc(func(_ context.Context, ref Ref) Blob {
		return FromFS(fsys, path.Clean(strings.Join(ref, "/")))
	})
}
