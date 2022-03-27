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
)

func Test(t *testing.T, s blob.Storage) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// TODO: actually test things...

	// TODO: log errors
	s.Set(ctx, "a", strings.NewReader("av"))
	s.Set(ctx, "b", strings.NewReader("bv"))
	s.Set(ctx, "d", strings.NewReader("dv"))
	s.Set(ctx, "c", strings.NewReader("cv"))
	s.Set(ctx, "e", strings.NewReader("ev"))
	s.Set(ctx, "f", strings.NewReader("fv"))
	s.Set(ctx, "g", strings.NewReader("gv"))
	s.Set(ctx, "h", strings.NewReader("hv"))

	s.Delete(ctx, "b")
	t.Log(blobStr(s.Get(ctx, "b")))

	// TODO: test stat
	testEnumerate(ctx, t, s)
}

func testEnumerate(ctx context.Context, t *testing.T, s blob.Storage) {
	if err := s.Enumerate(context.Background(), "", func(ref string, size int64) error {
		header, err := blob.Peek(s.Get(ctx, ref), 1<<10)
		if err != nil {
			return fmt.Errorf("%v: %v", ref, err)
		}
		t.Log(ref, size, http.DetectContentType(header))
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
