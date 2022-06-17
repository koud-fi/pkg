package search_test

import (
	"testing"

	"github.com/koud-fi/pkg/search"
)

func TestTagIndex(t *testing.T) {
	idx := search.NewMemoryTagIndex()
	idx.Put(search.Entry{ID: "1", Tags: []string{"a"}})
	idx.Put(search.Entry{ID: "4", Tags: []string{"b", "d"}})
	idx.Put(search.Entry{ID: "2", Tags: []string{"a", "b", "c"}})
	idx.Put(search.Entry{ID: "3", Tags: []string{"b", "c"}})

	t.Log(idx.Query([]string{"a"}, 10))
	t.Log(idx.Query([]string{"b"}, 10))
}
