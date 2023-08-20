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
	mu sync.RWMutex

	// TODO: optimize storage format to make insert/delete faster, use blob.Ref as pair key

	data refDataSlice
}

var _ blob.SortedStorage = (*Storage)(nil)

func NewStorage() *Storage {
	return &Storage{data: make(refDataSlice, 0)}
}

func (s *Storage) Get(_ context.Context, ref string) blob.Blob {
	return blob.ByteFunc(func() ([]byte, error) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		if i, ok := s.search(ref, true); ok {
			return s.data[i].Value(), nil
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

	p := rx.NewPair(ref, data)
	if i, ok := s.search(ref, false); ok {
		s.data[i] = p
	} else {
		s.data = append(s.data[:i], append(refDataSlice{p}, s.data[i:]...)...)
	}
	return nil
}

func (s *Storage) Delete(_ context.Context, refs ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ref := range refs {
		if i, ok := s.search(ref, true); ok {
			s.data = append(s.data[:i], s.data[i+1:]...)
		}
	}
	return nil
}

func (s *Storage) Iter(ctx context.Context, after string) rx.Iter[blob.RefBlob] {
	var i = -1

	return rx.FuncIter(func(rx.Done) ([]blob.RefBlob, rx.Done, error) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		if i < 0 {
			i, _ = s.search(after, true)
		}
		select {
		case <-ctx.Done():
			return nil, true, ctx.Err()
		default:
			var out []blob.RefBlob // TODO: return larger batches of data
			if i < len(s.data) {
				out = append(out, blob.RefBlob{
					Ref:  s.data[i].Key(),
					Blob: blob.FromBytes(s.data[i].Value()),
				})
				i++
			}
			return out, i == len(s.data), nil
		}
	})
}

func (s *Storage) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = nil
}

func (s *Storage) search(ref string, readonly bool) (int, bool) {
	i := sort.Search(len(s.data), func(i int) bool {
		return s.data[i].Key() >= ref
	})
	return i, i < len(s.data) && (s.data)[i].Key() == ref
}

type refDataSlice []rx.Pair[string, []byte]

func (s refDataSlice) Len() int           { return len(s) }
func (s refDataSlice) Less(i, j int) bool { return s[j].Key() < s[i].Key() }
func (s refDataSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
