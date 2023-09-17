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

func (s *Storage) Iter(ctx context.Context, state rx.Lens[string]) rx.Iter[blob.RefBlob] {
	return &iter{s: s, ctx: ctx, state: state, i: -1}
}

type iter struct {
	s     *Storage
	ctx   context.Context
	state rx.Lens[string]

	init bool
	i    int
	err  error
}

func (it *iter) Next() bool {
	it.s.mu.RLock()
	defer it.s.mu.RUnlock()

	select {
	case <-it.ctx.Done():
		it.err = it.ctx.Err()
		return false
	default:
		if !it.init {
			after, err := it.state.Get()
			if err != nil {
				it.err = err
				return false
			}
			it.i, _ = it.s.search(after, true)
			it.init = true
		} else {
			if it.i < len(it.s.data) {
				if it.err = it.state.Set(it.s.data[it.i].Key()); it.err == nil {
					it.i++
				}
			}
		}
		return it.err == nil && it.i < len(it.s.data)
	}
}

func (it *iter) Value() blob.RefBlob {
	it.s.mu.RLock()
	defer it.s.mu.RUnlock()

	if it.i >= len(it.s.data) {
		return blob.RefBlob{Blob: blob.Empty()}
	}
	p := it.s.data[it.i]
	return blob.RefBlob{
		Ref:  p.Key(),
		Blob: blob.FromBytes(p.Value()),
	}
}

func (it *iter) Close() error { return it.err }

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
