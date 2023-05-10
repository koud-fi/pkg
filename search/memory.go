package search

import (
	"sort"
	"strings"
	"sync"

	"github.com/koud-fi/pkg/bloom"
	"github.com/koud-fi/pkg/sorted"
)

const bloomFilterK = 9

type memTagIdx struct {
	mu           sync.RWMutex
	data         []memEntry
	dataIndex    map[string]int
	tagIDs       map[string]uint32
	tagCounts    map[string]int // TODO: more optimal tag store (prefix map?)
	isDirty      bool
	orderCounter int64
}

type memEntry struct {
	Entry     // TODO: avoid storing the full entry with all the tags etc.
	tagIDs    *sorted.Set[uint32]
	bloom     bloom.Filter
	isDeleted bool
}

func NewMemoryTagIndex() TagIndex {
	return &memTagIdx{
		dataIndex: make(map[string]int, 1<<8),
		tagIDs:    make(map[string]uint32, 1<<8),
		tagCounts: make(map[string]int, 1<<8),
	}
}

func (mti *memTagIdx) Query(tags []string, limit int) (QueryResult, error) {
	mti.Commit()

	mti.mu.RLock()
	defer mti.mu.RUnlock()

	qTagIDs, ok := mti.resolveTagIDs(tags, false)
	if !ok {
		return QueryResult{Data: []Entry{}}, nil
	}
	preAlloc := limit
	if preAlloc > 1<<10 {
		preAlloc = 1 << 10
	}
	var (
		qRes   = QueryResult{Data: make([]Entry, 0, preAlloc)}
		qBloom = bloom.New32(qTagIDs.Data(), bloomFilterK)
	)
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
		qRes.TotalCount++
		if limit > 0 && len(qRes.Data) == limit {
			continue
		}
		qRes.Data = append(qRes.Data, me.Entry)
	}
	return qRes, nil
}

func (mti *memTagIdx) Put(e ...Entry) {
	mti.mu.Lock()
	defer mti.mu.Unlock()

	for _, e := range e {

		// TODO: compare to existing entry to avoid pointless commits

		var (
			tagIDs, _ = mti.resolveTagIDs(e.Tags, true)
			me        = memEntry{
				Entry:  e,
				tagIDs: tagIDs,
				bloom:  bloom.New32(tagIDs.Data(), bloomFilterK),
			}
		)
		if i, ok := mti.dataIndex[e.ID]; ok {
			if me.Order <= 0 {
				me.isDeleted = me.Order < 0
				me.Order = mti.data[i].Order
			}

			// TODO: update tag counts

			mti.data[i] = me
		} else {
			if me.Order == 0 {
				mti.orderCounter++
				me.Order = -mti.orderCounter
			}
			for _, tag := range me.Tags {
				mti.tagCounts[tag]++
			}
			mti.dataIndex[e.ID] = len(mti.data)
			mti.data = append(mti.data, me)
		}
		mti.isDirty = true
	}
}

func (mti *memTagIdx) resolveTagIDs(tags []string, create bool) (*sorted.Set[uint32], bool) {
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

func (mti *memTagIdx) Commit() error {
	mti.mu.Lock()
	defer mti.mu.Unlock()

	if !mti.isDirty {
		return nil
	}
	sort.Slice(mti.data, func(i, j int) bool {
		if mti.data[i].Order == mti.data[j].Order {
			return mti.data[i].ID < mti.data[j].ID
		}
		return mti.data[i].Order > mti.data[j].Order
	})
	if len(mti.dataIndex) > len(mti.data)*2 {
		mti.dataIndex = make(map[string]int, len(mti.data))
	}
	for i, e := range mti.data {
		mti.dataIndex[e.ID] = i
	}
	mti.isDirty = false
	return nil
}

func (mti *memTagIdx) Tags(prefix string) ([]TagInfo, error) {
	var res []TagInfo // TODO: smart pre-alloc
	for tag, count := range mti.tagCounts {
		if strings.HasPrefix(tag, prefix) {
			res = append(res, TagInfo{Tag: tag, Count: count})
		}
	}
	return res, nil
}
