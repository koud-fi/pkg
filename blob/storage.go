package blob

import (
	"context"
	"io"
)

type Storage interface {
	Getter
	Setter
	Enumerator
	Statter
	Deleter
}

type Getter interface {
	Get(ctx context.Context, ref string) Blob
}

type Setter interface {
	Set(ctx context.Context, ref string, r io.Reader) error
}

type Enumerator interface {
	Enumerate(ctx context.Context, after string, fn func(string, int64) error) error
}

type Statter interface {
	Stat(ctx context.Context, refs []string, fn func(string, int64) error) error
}

type Deleter interface {
	Delete(ctx context.Context, refs ...string) error
}
