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
	"github.com/koud-fi/pkg/rx/diriter"
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

func DirPerm(m os.FileMode) Option          { return func(s *Storage) { s.dirPerm = m } }
func FilePerm(m os.FileMode) Option         { return func(s *Storage) { s.filePerm = m } }
func IterOpts(opt ...diriter.Option) Option { return func(s *Storage) { s.iterOpts = opt } }

var _ blob.SortedStorage = (*Storage)(nil)

type Storage struct {
	root            string
	bucketLevels    []int
	bucketPrefixLen int
	dirPerm         os.FileMode
	filePerm        os.FileMode
	iterOpts        []diriter.Option
}

func NewStorage(root string, opt ...Option) (*Storage, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	s := Storage{
		root:     absRoot,
		dirPerm:  defaultDirPerm,
		filePerm: defaultFilePerm,
	}
	for _, opt := range opt {
		opt(&s)
	}
	return &s, nil
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

func (s Storage) Iter(_ context.Context, after string) rx.Iter[blob.RefBlob] {
	if after != "" {
		panic("localdisk.Iter: after not supported") // TODO: implement "after"
	}
	d := diriter.New(os.DirFS(s.root), ".", s.iterOpts...)
	return rx.Map(d, (func(e diriter.Entry) blob.RefBlob {
		return blob.RefBlob{
			Ref:  e.Path,
			Blob: localfile.New(s.refPath(e.Path)),
		}
	}))
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
		parts := make([]string, 0, len(s.bucketLevels)+2)
		parts = append(parts, s.root)
		var start int
		for _, l := range s.bucketLevels {
			end := start + l
			if end > len(ref) {
				end = len(ref)
			}
			part := ref[start:end] + strings.Repeat("_", l-(end-start))
			parts = append(parts, part)
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
