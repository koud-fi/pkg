package localdisk

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/blob/localfile"
)

const (
	defaultDirPerm  = os.FileMode(0700)
	defaultFilePerm = os.FileMode(0600)
)

type Option func(*Storage)

func HideFunc(fn func(name string) bool) func(*Storage) { return func(s *Storage) { s.hideFunc = fn } }
func DirPerm(m os.FileMode) func(*Storage)              { return func(s *Storage) { s.dirPerm = m } }
func FilePerm(m os.FileMode) func(*Storage)             { return func(s *Storage) { s.filePerm = m } }

var _ blob.Storage = (*Storage)(nil)

type Storage struct {
	root     string
	hideFunc func(string) bool
	dirPerm  os.FileMode
	filePerm os.FileMode
}

func NewStorage(root string, opt ...Option) (*Storage, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	s := Storage{
		root:     absRoot,
		hideFunc: defaultHideFunc,
		dirPerm:  defaultDirPerm,
		filePerm: defaultFilePerm,
	}
	for _, opt := range opt {
		opt(&s)
	}
	return &s, nil
}

func defaultHideFunc(name string) bool {
	return strings.HasPrefix(name, ".")
}

func (s Storage) Fetch(_ context.Context, ref string) blob.Blob {
	return localfile.New(s.refPath(ref))
}

func (s Storage) Receive(_ context.Context, ref string, r io.Reader) error {
	path := s.refPath(ref)
	if err := os.MkdirAll(filepath.Dir(path), s.dirPerm); err != nil {
		return err
	}
	return localfile.WriteReader(path, r, s.filePerm)
}

func (s Storage) Enumerate(ctx context.Context, after string, fn func(string, int64) error) error {
	if after != "" {
		return errors.New("localdisk.Enumerate: after not supported") // TODO: implement "after"
	}
	return s.enumDir(ctx, s.root, fn)
}

func (s Storage) enumDir(ctx context.Context, dirPath string, fn func(string, int64) error) error {
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, d := range dir {
		if s.hideFunc(d.Name()) {
			continue
		}
		path := filepath.Join(dirPath, d.Name())
		if d.IsDir() {
			if err := s.enumDir(ctx, path, fn); err != nil {
				return err
			}
			continue
		}
		ref, err := s.pathRef(path)
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			stat, err := d.Info()
			if err != nil {
				return err
			}
			if err := fn(ref, stat.Size()); err != nil {
				return err
			}
		}
	}
	return nil
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

func (s Storage) pathRef(path string) (string, error) {
	return filepath.Rel(s.root, path)
}
