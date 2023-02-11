package localdisk

import (
	"context"
	"crypto"
	"encoding/hex"
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

func BucketHash(h crypto.Hash) Option { return func(s *Storage) { s.bucketHash = &h } }
func DirPerm(m os.FileMode) Option    { return func(s *Storage) { s.dirPerm = m } }
func FilePerm(m os.FileMode) Option   { return func(s *Storage) { s.filePerm = m } }

var _ blob.SortedStorage = (*Storage)(nil)

type Storage struct {
	root            string
	bucketLevels    []int
	bucketPrefixLen int
	bucketHash      *crypto.Hash
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
		dirPerm:  defaultDirPerm,
		filePerm: defaultFilePerm,
	}
	for _, opt := range opt {
		opt(&s)
	}
	return &s, nil
}

func (s Storage) Get(_ context.Context, ref blob.Ref) blob.Blob {
	return localfile.New(s.refPath(ref))
}

func (s Storage) Set(_ context.Context, ref blob.Ref, r io.Reader) error {
	path := s.refPath(ref)
	if err := os.MkdirAll(filepath.Dir(path), s.dirPerm); err != nil {
		return err
	}
	return localfile.WriteReader(path, r, s.filePerm)
}

func (s Storage) Iter(_ context.Context, d blob.Domain, after blob.Ref) rx.Iter[blob.RefBlob] {
	if len(after) != 0 {
		panic("localdisk.Iter: after not supported") // TODO: implement "after"
	}
	var prefix string
	if d != blob.Default {
		prefix = string(d) + "/"
	}
	it := diriter.New(os.DirFS(s.root), string(d))
	return rx.Map(it, (func(e diriter.Entry) blob.RefBlob {
		ref := blob.NewRef(d, strings.TrimPrefix(e.Path(), prefix))
		return blob.RefBlob{
			Ref:  ref,
			Blob: localfile.New(filepath.Join(s.root, ref.String())),
		}
	}))
}

func (s Storage) Delete(_ context.Context, refs ...blob.Ref) error {
	for _, ref := range refs {
		if err := os.Remove(s.refPath(ref)); err != nil {
			return err
		}
	}
	return nil
}

func (s Storage) refPath(ref blob.Ref) string {
	refStr := ref.Ref().String()
	if len(s.bucketLevels) > 0 {
		parts := make([]string, 0, len(s.bucketLevels)+2)
		parts = append(parts, s.root, string(ref.Domain()))
		var (
			bucketRef = refStr
			start     int
		)
		if s.bucketHash != nil {
			h := s.bucketHash.New()
			if _, err := h.Write([]byte(bucketRef)); err != nil {
				panic("localdisk.refPath bucket hash error: " + err.Error())
			}
			bucketRef = hex.EncodeToString(h.Sum(nil))
		}
		for _, l := range s.bucketLevels {
			end := start + l
			if end > len(bucketRef) {
				end = len(bucketRef)
			}
			part := bucketRef[start:end] + strings.Repeat("_", l-(end-start))
			parts = append(parts, part)
			start = end
		}
		parts = append(parts, refStr)
		return filepath.Join(parts...)
	}
	return filepath.Join(s.root, string(ref.Domain()), refStr)
}
