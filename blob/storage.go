package blob

import (
	"context"
	"io"
)

type Storage interface {
	Fetcher
	Receiver
	Enumerator
	Statter
	Remover
}

type Fetcher interface {
	Fetch(ctx context.Context, ref string) Blob
}

type Receiver interface {
	Receive(ctx context.Context, ref string, r io.Reader) error
}

type Enumerator interface {
	Enumerate(ctx context.Context, after string, fn func(string, int64) error) error
}

type Statter interface {
	Stat(ctx context.Context, refs []string, fn func(string, int64) error) error
}

type Remover interface {
	Remove(ctx context.Context, refs ...string) error
}
