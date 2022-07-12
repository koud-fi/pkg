package blob_test

import (
	"testing"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/rx"
)

func TestLineIter(t *testing.T) {
	it := blob.LineIter(blob.FromString("a\nb\nc"))
	if err := rx.ForEach(it, func(s string) error {
		t.Log(s)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
