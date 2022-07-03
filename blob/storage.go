package blob

import (
	"context"
	"io"

	"github.com/koud-fi/pkg/rx"
)

type Storage interface {
	Getter
	Setter
	Iterator
	Deleter
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
	Iter(ctx context.Context, after string) rx.Iter[RefBlob]
}
