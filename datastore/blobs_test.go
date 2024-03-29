package datastore_test

import (
	"testing"

	"github.com/koud-fi/pkg/datastore"
	"github.com/koud-fi/pkg/rx"
	"github.com/koud-fi/pkg/storage/memory"

	"golang.org/x/net/context"
)

type TestData struct {
	ID    string
	Value int
}

func TestBlobsTable(t *testing.T) {
	var (
		ctx = context.Background()
		bt  = datastore.BlobsTable(
			memory.NewStorage(),
			func(v TestData) (string, error) { return v.ID, nil })
	)
	t.Log(rx.Slice(rx.MapErr(rx.SliceIter(
		TestData{ID: "1", Value: 42},
		TestData{ID: "2", Value: 69}), bt.Put(ctx))))

	t.Log(bt.Delete(ctx)(TestData{ID: "1"}))

	t.Log(rx.Slice(rx.MapErr(rx.SliceIter(
		TestData{ID: "1"},
		TestData{ID: "2"},
		TestData{ID: "3"}), bt.Get(ctx))))
}
