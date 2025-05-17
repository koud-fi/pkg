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

func TestTagIndex(t *testing.T) {

	// TODO: proper test with actual assertion

	adapter := search.NewAdapter(
		func(te testEntry) string { return te.id },
		func(te testEntry) []string { return te.tags },
	)
	idx := search.NewShardedTagIndex(
		adapter,
		32,
		func(adapter search.Adapter[testEntry], _ int32) search.TagIndex[testEntry] {
			return search.NewMemoryTagIndex(adapter)
		},
	)
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
