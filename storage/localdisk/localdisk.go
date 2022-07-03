package localdisk

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/blob/localfile"
	"github.com/koud-fi/pkg/rx"
)

const (
	defaultDirPerm  = os.FileMode(0700)
	defaultFilePerm = os.FileMode(0600)
)

type Option func(*Storage)

func Buckets(levels ...int) Option {
	return func(s *Storage) {
		s.bucketLevels = levels
		s.bucketPrefixLen = 0
		for _, l := range levels {
			s.bucketPrefixLen += l + 1
		}
	}
}

func HideFunc(fn func(name string) bool) Option { return func(s *Storage) { s.hideFunc = fn } }
func DirPerm(m os.FileMode) Option              { return func(s *Storage) { s.dirPerm = m } }
func FilePerm(m os.FileMode) Option             { return func(s *Storage) { s.filePerm = m } }

var _ blob.Storage = (*Storage)(nil)

type Storage struct {
	root            string
	bucketLevels    []int
	bucketPrefixLen int
	hideFunc        func(string) bool
	dirPerm         os.FileMode
	filePerm        os.FileMode
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

func (s Storage) Get(_ context.Context, ref string) blob.Blob {
	return localfile.New(s.refPath(ref))
}

func (s Storage) Set(_ context.Context, ref string, r io.Reader) error {
	path := s.refPath(ref)
	if err := os.MkdirAll(filepath.Dir(path), s.dirPerm); err != nil {
		return err
	}
	return localfile.WriteReader(path, r, s.filePerm)
}

/*
func (s Storage) Enumerate(ctx context.Context, after string, fn func(string, int64) error) error {
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
*/

func (s *Storage) Iter(ctx context.Context, after string) rx.Iter[blob.RefBlob] {
	if after != "" {
		panic("localdisk.Iter: after not supported") // TODO: implement "after"
	}

	// ???

	panic("TODO")
}

func (s Storage) Delete(_ context.Context, refs ...string) error {
	for _, ref := range refs {
		if err := os.Remove(s.refPath(ref)); err != nil {
			return err
		}
	}
	return nil
}

func (s Storage) refPath(ref string) string {
	if len(s.bucketLevels) > 0 {
		parts := make([]string, len(s.bucketLevels)+2)
		parts = append(parts, s.root)
		var start int
		for _, l := range s.bucketLevels {
			end := start + l
			if end > len(ref) {
				end = len(ref)
			}
			parts = append(parts, ref[start:end]+strings.Repeat("_", l-(end-start)))
			start = end
		}
		parts = append(parts, ref)
		return filepath.Join(parts...)
	}
	return filepath.Join(s.root, ref)
}

func (s Storage) pathRef(path string) (string, error) {
	ref, err := filepath.Rel(s.root, path)
	if err != nil {
		return "", err
	}
	return ref[s.bucketPrefixLen:], nil
}
