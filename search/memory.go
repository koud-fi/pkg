package search

import (
	"sort"
	"strings"
	"sync"

	"github.com/koud-fi/pkg/bloom"
	"github.com/koud-fi/pkg/sorted"
)

const bloomFilterK = 9

type memTagIdx[T Entry] struct {
	mu           sync.RWMutex
	data         []memEntry[T]
	dataIndex    map[string]int
	tagIDs       map[string]uint32
	tagCounts    map[string]int // TODO: more optimal tag store (prefix map?)
	isDirty      bool
	orderCounter int64
}

type memEntry[T any] struct {
	entry     T
	tagIDs    *sorted.Set[uint32]
	order     int64
	bloom     bloom.Filter
	isDeleted bool
}

func NewMemoryTagIndex[T Entry]() TagIndex[T] {
	return &memTagIdx[T]{
		dataIndex: make(map[string]int, 1<<8),
		tagIDs:    make(map[string]uint32, 1<<8),
		tagCounts: make(map[string]int, 1<<8),
	}
}

func (mti *memTagIdx[T]) Get(id ...string) ([]T, error) {
	mti.mu.RLock()
	defer mti.mu.RUnlock()

	out := make([]T, 0, len(id))
	for _, id := range id {
		if idx, ok := mti.dataIndex[id]; ok {
			out = append(out, mti.data[idx].entry)
		}
	}
	return out, nil
}

func (mti *memTagIdx[T]) Query(dst *QueryResult[T], tags []string, limit int) error {
	mti.mu.RLock()
	defer mti.mu.RUnlock()

	qTagIDs, ok := mti.resolveTagIDs(tags, false)
	if !ok {
		return nil
	}
	dst.Reset()
	limit = max(0, limit)

	qBloom := bloom.New32(qTagIDs.Data(), bloomFilterK)
	for _, me := range mti.data {
		if me.isDeleted {
			continue
		}
		if !me.bloom.Contains(qBloom) {
			continue
		}
		if !me.tagIDs.HasSubset(qTagIDs) {
			continue
		}
		dst.TotalCount++
		if len(dst.Data) == limit {
			continue
		}
		dst.Data = append(dst.Data, me.entry)
	}
	return nil
}

func (mti *memTagIdx[T]) Put(e ...T) {
	mti.mu.Lock()
	defer mti.mu.Unlock()

	for _, e := range e {

		// TODO: compare to existing entry to avoid pointless commits

		var (
			tagIDs, _ = mti.resolveTagIDs(e.Tags(), true)
			me        = memEntry[T]{
				entry:  e,
				tagIDs: tagIDs,
				bloom:  bloom.New32(tagIDs.Data(), bloomFilterK),
			}
		)
		if ord, ok := any(e).(OrderedEntry); ok {
			me.order = ord.Order()
		}
		if i, ok := mti.dataIndex[e.ID()]; ok {
			if me.order <= 0 {
				me.isDeleted = me.order < 0
				me.order = mti.data[i].order
			}

			// TODO: update tag counts

			mti.data[i] = me
		} else {
			if ord, ok := any(e).(OrderedEntry); ok {
				me.order = ord.Order()
			}
			if me.order == 0 {
				mti.orderCounter++
				me.order = -mti.orderCounter
			}
			for _, tag := range e.Tags() {
				mti.tagCounts[tag]++
			}
			mti.dataIndex[e.ID()] = len(mti.data)
			mti.data = append(mti.data, me)
		}
		mti.isDirty = true
	}
}

func (mti *memTagIdx[_]) resolveTagIDs(tags []string, create bool) (*sorted.Set[uint32], bool) {
	ids := make([]uint32, 0, len(tags))
	for _, tag := range tags {
		if len(tag) == 0 {
			continue
		}
		id, ok := mti.tagIDs[tag]
		if !ok {
			if !create {
				return nil, false
			}
			id = uint32(len(mti.tagIDs) + 1)
			mti.tagIDs[tag] = id
		}
		ids = append(ids, id)
	}
	return sorted.NewSet(ids...), true
}

func (mti *memTagIdx[_]) Commit() error {
	mti.mu.Lock()
	defer mti.mu.Unlock()

	if !mti.isDirty {
		return nil
	}
	sort.Slice(mti.data, func(i, j int) bool {
		if mti.data[i].order == mti.data[j].order {
			return mti.data[i].entry.ID() < mti.data[j].entry.ID()
		}
		return mti.data[i].order > mti.data[j].order
	})
	if len(mti.dataIndex) > len(mti.data)*2 {
		mti.dataIndex = make(map[string]int, len(mti.data))
	}
	for i, e := range mti.data {
		mti.dataIndex[e.entry.ID()] = i
	}
	mti.isDirty = false
	return nil
}

func (mti *memTagIdx[_]) Tags(prefix string) ([]TagInfo, error) {
	var res []TagInfo // TODO: smart pre-alloc
	for tag, count := range mti.tagCounts {
		if strings.HasPrefix(tag, prefix) {
			res = append(res, TagInfo{Tag: tag, Count: count})
		}
	}
	return res, nil
}
