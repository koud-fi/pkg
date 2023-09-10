package storagetest

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/datastore"
	"github.com/koud-fi/pkg/rx"
)

// TODO: actually test things... (needs some sort of assertion)

func Test(t *testing.T, s blob.Storage) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	testKV(t, ctx, datastore.BlobsKV[string](s))
	testIter(ctx, t, s)
}

func testIter(ctx context.Context, t *testing.T, s blob.Storage) {
	ss, ok := s.(blob.SortedStorage)
	if !ok {
		return
	}
	if err := rx.ForEach(ss.Iter(ctx, ""), func(b blob.RefBlob) error {
		header, err := blob.Peek(s.Get(ctx, b.Ref), 1<<10)
		if err != nil {
			return fmt.Errorf("%v: %v", b.Ref, err)
		}
		t.Log(b.Ref, http.DetectContentType(header))
		return nil

	}); err != nil {
		t.Fatal(err)
	}
}

func TestKV(t *testing.T, kv datastore.KV[string]) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	testKV(t, ctx, kv)
}

func testKV(t *testing.T, ctx context.Context, kv datastore.KV[string]) {

	// TODO: don't log nil errors

	t.Log(kv.Put(ctx, "a", "av"))
	t.Log(kv.Put(ctx, "bb", "bv"))
	t.Log(kv.Put(ctx, "dddd", "dv"))
	t.Log(kv.Put(ctx, "ccc", "cv"))
	t.Log(kv.Put(ctx, "eeeee", "ev"))
	t.Log(kv.Put(ctx, "ffffff", "fv"))
	t.Log(kv.Put(ctx, "gggg/ggg", "gv"))
	t.Log(kv.Put(ctx, "hhhh/hhhh", "hv"))

	t.Log(kv.Delete(ctx, "bb"))

	t.Log(kv.Get(ctx, "a"))
	t.Log(kv.Get(ctx, "bb"))
	t.Log(kv.Get(ctx, "gggg/ggg"))

	// TODO: iterator test (deprecates testIter)
}
