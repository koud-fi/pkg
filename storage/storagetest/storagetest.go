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
	"github.com/koud-fi/pkg/rx/lens"
)

// TODO: actually test things... (needs some sort of assertion)

func Test(t *testing.T, s blob.Storage) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	testKV(t, ctx, datastore.BlobsKV[string](s))
	testIter(t, ctx, s)
}

func testIter(t *testing.T, ctx context.Context, s blob.Storage) {
	ss, ok := s.(blob.SortedStorage)
	if !ok {
		return
	}
	state := lens.Value("")
	if err := rx.ForEach(ss.Iter(ctx, state), func(b blob.RefBlob) error {
		header, err := blob.Peek(ss.Get(ctx, b.Ref), 1<<10)
		if err != nil {
			return fmt.Errorf("%v: %v", b.Ref, err)
		}
		t.Log(b.Ref, http.DetectContentType(header))
		return nil

	}); err != nil {
		t.Fatal(err)
	}
	t.Log(state.Get())
}

func TestKV(t *testing.T, kv datastore.KV[string]) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	testKV(t, ctx, kv)
	testKVIter(t, ctx, kv)
}

func testKVIter(t *testing.T, ctx context.Context, kv datastore.KV[string]) {
	skv, ok := kv.(datastore.SortedKV[string])
	if !ok {
		return
	}
	state := lens.Value("")
	if err := rx.ForEach(skv.Iter(ctx, state), func(p rx.Pair[string, string]) error {
		t.Log(p.Key(), p.Value())
		return nil

	}); err != nil {
		t.Fatal(err)
	}
	t.Log(state.Get())
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
