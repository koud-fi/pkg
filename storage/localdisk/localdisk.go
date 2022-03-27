package localdisk

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/blob/localfile"
)

const (
	defaultDirPerm  = os.FileMode(0700)
	defaultFilePerm = os.FileMode(0600)
)

var _ blob.Storage = (*Storage)(nil)

type Storage struct {
	root     string
	dirPerm  os.FileMode
	filePerm os.FileMode
}

func NewStorage(root string) (*Storage, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	return &Storage{
		root:     absRoot,
		dirPerm:  defaultDirPerm,
		filePerm: defaultFilePerm,
	}, nil
}

func (s Storage) Fetch(_ context.Context, ref string) blob.Blob {
	return localfile.New(s.refPath(ref))
}

func (s Storage) Receive(_ context.Context, ref string, r io.Reader) error {
	return localfile.WriteReader(s.refPath(ref), r, s.filePerm)
}

func (s Storage) Enumerate(ctx context.Context, after string, fn func(string, int64) error) error {

	// ???

	panic("TODO")
}

func (s Storage) Stat(_ context.Context, refs []string, fn func(string, int64) error) error {
	for _, ref := range refs {
		info, err := os.Stat(s.refPath(ref))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		if err := fn(ref, info.Size()); err != nil {
			return err
		}
	}
	return nil
}

func (s Storage) Remove(_ context.Context, refs ...string) error {
	for _, ref := range refs {
		if err := os.Remove(s.refPath(ref)); err != nil {
			return err
		}
	}
	return nil
}

func (s Storage) refPath(ref string) string {
	return filepath.Join(s.root, ref)
}
