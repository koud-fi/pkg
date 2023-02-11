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

	data map[blob.Domain]*refDataSlice
}

var _ blob.SortedStorage = (*Storage)(nil)

func NewStorage() *Storage {
	return &Storage{data: make(map[blob.Domain]*refDataSlice)}
}

func (s *Storage) Get(_ context.Context, ref blob.Ref) blob.Blob {
	return blob.ByteFunc(func() ([]byte, error) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		if i, data, ok := s.search(ref, true); ok {
			return (*data)[i].Value(), nil
		}
		return nil, os.ErrNotExist
	})
}

func (s *Storage) Set(_ context.Context, ref blob.Ref, r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	p := rx.NewPair(ref.String(), data)
	if i, data, ok := s.search(ref, false); ok {
		(*data)[i] = p
	} else {
		(*data) = append((*data)[:i], append(refDataSlice{p}, (*data)[i:]...)...)
	}
	return nil
}

func (s *Storage) Iter(ctx context.Context, d blob.Domain, after blob.Ref) rx.Iter[blob.RefBlob] {
	var (
		i    = -1
		data *refDataSlice
	)
	return rx.FuncIter(func(rx.Done) ([]blob.RefBlob, rx.Done, error) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		if i < 0 {
			i, data, _ = s.search(blob.NewRef(d, after.Ref()...), true)
		}
		if data == nil {
			return nil, true, nil
		}
		select {
		case <-ctx.Done():
			return nil, true, ctx.Err()
		default:
			var out []blob.RefBlob // TODO: return larger batches of data
			if i < len(*data) {
				out = append(out, blob.RefBlob{
					Ref:  blob.ParseRef((*data)[i].Key()),
					Blob: blob.FromBytes((*data)[i].Value()),
				})
				i++
			}
			return out, i == len(*data), nil
		}
	})
}

func (s *Storage) Delete(_ context.Context, refs ...blob.Ref) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ref := range refs {
		if i, data, ok := s.search(ref, true); ok {
			(*data) = append((*data)[:i], (*data)[i+1:]...)
		}
	}
	return nil
}

func (s *Storage) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = nil
}

func (s *Storage) search(ref blob.Ref, readonly bool) (int, *refDataSlice, bool) {
	d := ref.Domain()
	if d == "" {
		d = blob.Domain(ref.String())
	}
	data, ok := s.data[d]
	if !ok {
		if readonly {
			return 0, nil, false
		}
		data = &refDataSlice{}
		s.data[ref.Domain()] = data
		return 0, data, false
	}
	refStr := ref.String()
	i := sort.Search(len(*data), func(i int) bool {
		return (*data)[i].Key() >= refStr
	})
	return i, data, i < len(*data) && (*data)[i].Key() == refStr
}

type refDataSlice []rx.Pair[string, []byte]

func (s refDataSlice) Len() int           { return len(s) }
func (s refDataSlice) Less(i, j int) bool { return s[j].Key() < s[i].Key() }
func (s refDataSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
