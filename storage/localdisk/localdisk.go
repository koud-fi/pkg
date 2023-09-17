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

func (s Storage) Delete(_ context.Context, refs ...string) error {
	for _, ref := range refs {
		if err := os.Remove(s.refPath(ref)); err != nil {
			return err
		}
	}
	return nil
}

func (s Storage) Iter(_ context.Context, state rx.Lens[string]) rx.Iter[blob.RefBlob] {

	// TODO: implement proper state logic

	it := diriter.New(os.DirFS(s.root), ".")
	return rx.Map(it, (func(e diriter.Entry) blob.RefBlob {
		var (
			p  = e.Key()
			bl = len(s.bucketLevels)
		)
		if bl > 0 {
			p = strings.SplitN(p, "/", bl+1)[bl]
		}
		return blob.RefBlob{
			Ref:  p,
			Blob: localfile.New(filepath.Join(s.root, p)),
		}
	}))
}

func (s Storage) refPath(ref string) string {
	if len(s.bucketLevels) > 0 {
		parts := make([]string, 0, len(s.bucketLevels)+2)
		parts = append(parts, s.root)
		var (
			bucketRef = ref
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
		parts = append(parts, ref)
		return filepath.Join(parts...)
	}
	return filepath.Join(s.root, ref)
}
