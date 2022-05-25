package search

import "sync"

type memTagIdx struct {
	mu           sync.RWMutex
	data         []Entry
	dataIndex    map[string]int
	tagIDs       map[string]uint32
	tagCounts    map[string]int // TODO: more optimal tag store (prefix map?)
	isDirty      bool
	orderCounter int64
}

func NewMemoryTagIndex() TagIndex {
	return &memTagIdx{
		dataIndex: make(map[string]int, 1<<8),
		tagIDs:    make(map[string]uint32, 1<<8),
		tagCounts: make(map[string]int, 1<<8),
	}
}

func (mti *memTagIdx) Query(tags []string, limit int) (QueryResult, error) {

	// ???

	panic("TODO")
}

func (mti *memTagIdx) Put(e ...Entry) error {

	// ???

	panic("TODO")
}

func (mti *memTagIdx) Commit() error {

	// ???

	panic("TODO")
}

func (mti *memTagIdx) Tags(prefix string, limit int) ([]TagInfo, error) {

	// ???

	panic("TODO")
}
