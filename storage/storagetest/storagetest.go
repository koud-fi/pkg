package storagetest

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/rx"
)

func Test(t *testing.T, s blob.Storage) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// TODO: actually test things...

	// TODO: log errors
	s.Set(ctx, "a", strings.NewReader("av"))
	s.Set(ctx, "bb", strings.NewReader("bv"))
	s.Set(ctx, "dddd", strings.NewReader("dv"))
	s.Set(ctx, "ccc", strings.NewReader("cv"))
	s.Set(ctx, "eeeee", strings.NewReader("ev"))
	s.Set(ctx, "ffffff", strings.NewReader("fv"))
	s.Set(ctx, "gggg/ggg", strings.NewReader("gv"))
	s.Set(ctx, "hhhh/hhhh", strings.NewReader("hv"))

	s.Delete(ctx, "bb")
	t.Log(blobStr(s.Get(ctx, "a")))
	t.Log(blobStr(s.Get(ctx, "bb")))
	t.Log(blobStr(s.Get(ctx, "gggg/ggg")))

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

func blobStr(b blob.Blob) string {
	data, err := blob.Peek(b, 50)
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	for i := range data {
		if unicode.IsControl(rune(data[i])) {
			data[i] = ' '
		}
	}
	return string(data)
}
