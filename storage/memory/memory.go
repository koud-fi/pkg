package memory

import (
	"context"
	"io"
	"os"
	"sort"
	"sync"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/rx"
)

type Storage struct {
	mu   sync.RWMutex
	data refDataSlice // TODO: optimize storage format to make insert/delete faster
}

var _ blob.Storage = (*Storage)(nil)

func NewStorage() *Storage {
	return new(Storage)
}

func (s *Storage) Get(_ context.Context, ref string) blob.Blob {
	return blob.ByteFunc(func() ([]byte, error) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		if i, ok := s.search(ref); ok {
			return s.data[i].data, nil
		}
		return nil, os.ErrNotExist
	})
}

func (s *Storage) Set(_ context.Context, ref string, r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	if i, ok := s.search(ref); ok {
		s.data[i].data = data
	} else {
		v := refData{ref: ref, data: data}
		s.data = append(s.data[:i], append(refDataSlice{v}, s.data[i:]...)...)
	}
	return nil
}

func (s *Storage) Iter(ctx context.Context, after string) rx.Iter[blob.RefBlob] {
	i := -1
	return rx.FuncIter(func() ([]blob.RefBlob, bool, error) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		if i < 0 {
			i, _ = s.search(after)
		}
		select {
		case <-ctx.Done():
			return nil, false, ctx.Err()
		default:
			var out []blob.RefBlob // TODO: return larger batches of data
			if i < len(s.data) {
				out = append(out, blob.RefBlob{
					Ref:  s.data[i].ref,
					Blob: blob.FromBytes(s.data[i].data),
				})
				i++
			}
			return out, i < len(s.data), nil
		}
	})
}

func (s *Storage) Delete(_ context.Context, refs ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ref := range refs {
		if i, ok := s.search(ref); ok {
			s.data = append(s.data[:i], s.data[i+1:]...)
		}
	}
	return nil
}

func (s *Storage) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = s.data[:0]
}

func (s *Storage) search(ref string) (int, bool) {
	i := sort.Search(len(s.data), func(i int) bool {
		return s.data[i].ref >= ref
	})
	return i, i < len(s.data) && s.data[i].ref == ref
}

type refData struct {
	ref  string
	data []byte
}

type refDataSlice []refData

func (s refDataSlice) Len() int           { return len(s) }
func (s refDataSlice) Less(i, j int) bool { return s[j].ref < s[i].ref }
func (s refDataSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
