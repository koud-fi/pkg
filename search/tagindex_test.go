package search_test

import (
	"testing"

	"github.com/koud-fi/pkg/search"
)

type testEntry struct {
	id    string
	tags  []string
	order int64
}

func (te testEntry) ID() string     { return te.id }
func (te testEntry) Tags() []string { return te.tags }
func (te testEntry) Order() int64   { return te.order }

func TestTagIndex(t *testing.T) {

	// TODO: proper test with actual assertion

	idx := search.NewShardedTagIndex[testEntry](32, func(_ int32) search.TagIndex[testEntry] {
		return search.NewMemoryTagIndex[testEntry]()
	})
	idx.Put(testEntry{"1", []string{"a"}, 4})
	idx.Put(testEntry{"4", []string{"b", "d"}, 1})
	idx.Put(testEntry{"2", []string{"a", "b", "c"}, 3})
	idx.Put(testEntry{"3", []string{"b", "c"}, 2})

	var res search.QueryResult[testEntry]
	t.Log(idx.Query(&res, []string{"a"}, 10))
	t.Log(res)
	t.Log(idx.Query(&res, []string{"b"}, 10))
	t.Log(res)

	t.Log(idx.Get("0", "1", "2"))
}
