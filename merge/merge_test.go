package merge_test

import (
	"net/url"
	"testing"

	"github.com/koud-fi/pkg/merge"
)

type dstType struct {
	A int
	B string
}

func TestMerge(t *testing.T) {
	var dst dstType
	if err := merge.Values(&dst, url.Values{
		"A": {"42"},
		"B": {"Hello, world?"},
	}.Get); err != nil {
		t.Fatal(err)
	}
	t.Log(dst) // TODO: assert the result correctly
}
