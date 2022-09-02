package redis

import (
	"context"
	"io"

	"github.com/koud-fi/pkg/blob"
)

var _ blob.Storage = (*Storage)(nil)

type Storage struct {
	// TODO
}

func (s *Storage) Get(ctx context.Context, ref string) blob.Blob {
	panic("TODO")
}

func (s *Storage) Set(ctx context.Context, ref string, r io.Reader) error {
	panic("TODO")
}

func (s *Storage) Delete(ctx context.Context, refs ...string) error {
	panic("TODO")
}
