package kv

import (
	"container/list"
	"context"
	"iter"
	"sync"

	"github.com/koud-fi/pkg/errx"
)

var _ interface {
	Storage[string, any]
	ScanReader[string, any]
} = (*MemoryStorage[string, any])(nil)

type (
	MemoryStorage[K comparable, V any] struct {
		kind string

		mu             sync.RWMutex
		dataList       *list.List
		dataMap        map[K]*list.Element
		versionCounter uint64
	}
	memoryPair[K comparable, V any] struct {
		key     K
		version uint64
		value   V
	}
)

func (p memoryPair[K, V]) Key() K   { return p.key }
func (p memoryPair[K, V]) Value() V { return p.value }

func NewMemoryStorage[K comparable, V any](kind string) *MemoryStorage[K, V] {
	return &MemoryStorage[K, V]{
		dataList:       list.New(),
		dataMap:        make(map[K]*list.Element),
		versionCounter: 1,
	}
}

func (s *MemoryStorage[K, V]) Get(ctx context.Context, key K) (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	el, ok := s.dataMap[key]
	if !ok {
		var zero V
		return zero, errx.NewNotFound(s.kind, key)
	}
	return el.Value.(memoryPair[K, V]).value, nil
}

func (s *MemoryStorage[K, V]) Scan(ctx context.Context) (iter.Seq[Pair[K, V]], func() error) {

	// TODO: context should be able to cancel the running iterator
	// TODO: there may be race conditions here

	return func(yield func(Pair[K, V]) bool) {
		for el := s.dataList.Front(); el != nil; el = el.Next() {
			s.mu.RLock()
			p := el.Value.(memoryPair[K, V])
			s.mu.RUnlock()

			if !yield(p) {
				break
			}
		}
	}, func() error { return nil }
}

/*
func (s *MemoryStorage[K, V]) seek(afterVersion uint64) *list.Element {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if afterVersion == 0 {
		return s.dataList.Back()
	}

	// TODO: we should probably start from the front by default

	el := s.dataList.Front()
	if afterVersion == 0 {
		return el
	}
	for el != nil && el.Value.(memoryPair[K, V]).version <= afterVersion {
		el = el.Next()
	}
	return el
}
*/

func (s *MemoryStorage[K, V]) Set(_ context.Context, key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	newPair := memoryPair[K, V]{key, s.versionCounter, value}
	s.versionCounter++

	if el, ok := s.dataMap[newPair.key]; ok {
		el.Value = newPair
		s.dataList.MoveToFront(el)
		return nil
	}
	el := s.dataList.PushFront(newPair)
	s.dataMap[newPair.key] = el
	return nil
}

func (s *MemoryStorage[K, V]) Del(_ context.Context, key K) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	el, ok := s.dataMap[key]
	if !ok {
		return nil
	}
	s.dataList.Remove(el)
	delete(s.dataMap, key)
	return nil
}
