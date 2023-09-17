package blob

import (
	"context"
	"io"

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
	Get(ctx context.Context, ref string) Blob
}

type Setter interface {
	Set(ctx context.Context, ref string, r io.Reader) error
}

type Deleter interface {
	Delete(ctx context.Context, refs ...string) error
}

type RefBlob struct {
	Ref string
	Blob
}

type Iterator interface {
	Iter(ctx context.Context, state rx.Lens[string]) rx.Iter[RefBlob]
}
