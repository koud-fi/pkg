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

const testDomain = "__test"

func Test(t *testing.T, s blob.Storage) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// TODO: actually test things...

	// TODO: log errors
	s.Set(ctx, blob.NewRef(testDomain, "a"), strings.NewReader("av"))
	s.Set(ctx, blob.NewRef(testDomain, "bb"), strings.NewReader("bv"))
	s.Set(ctx, blob.NewRef(testDomain, "dddd"), strings.NewReader("dv"))
	s.Set(ctx, blob.NewRef(testDomain, "ccc"), strings.NewReader("cv"))
	s.Set(ctx, blob.NewRef(testDomain, "eeeee"), strings.NewReader("ev"))
	s.Set(ctx, blob.NewRef(testDomain, "ffffff"), strings.NewReader("fv"))
	s.Set(ctx, blob.NewRef(testDomain, "gggg/ggg"), strings.NewReader("gv"))
	s.Set(ctx, blob.NewRef(testDomain, "hhhh/hhhh"), strings.NewReader("hv"))

	s.Delete(ctx, blob.NewRef(testDomain, "bb"))
	t.Log(blobStr(s.Get(ctx, blob.NewRef(testDomain, "a"))))
	t.Log(blobStr(s.Get(ctx, blob.NewRef(testDomain, "bb"))))

	testIter(ctx, t, s)
}

func testIter(ctx context.Context, t *testing.T, s blob.Storage) {
	ss, ok := s.(blob.SortedStorage)
	if !ok {
		return
	}
	if err := rx.ForEach(ss.Iter(ctx, testDomain, nil), func(b blob.RefBlob) error {
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
