package localdisk

import (
	"context"
	"io"
	"path/filepath"

	"github.com/koud-fi/pkg/blob"
)

var _ blob.Storage = (*Storage)(nil)

type Storage struct {
	root string
}

func NewStorage(root string) (*Storage, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	return &Storage{root: absRoot}, nil
}

func (s Storage) Fetch(_ context.Context, ref string) blob.Blob {

	// ???

	panic("TODO")
}

func (s Storage) Receive(_ context.Context, ref string, r io.Reader) error {

	// ???

	panic("TODO")
}

func (s Storage) Enumerate(ctx context.Context, after string, fn func(string, int64) error) error {

	// ???

	panic("TODO")
}

func (s Storage) Stat(_ context.Context, refs []string, fn func(string, int64) error) error {

	// ???

	panic("TODO")
}

func (s Storage) Remove(_ context.Context, refs ...string) error {

	// ???

	panic("TODO")
}
