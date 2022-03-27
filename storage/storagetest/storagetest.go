package storagetest

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"unicode"

	"github.com/koud-fi/pkg/blob"
)

func Test(t *testing.T, s blob.Storage) {
	ctx := context.Background()

	// TODO: actually test things...

	// TODO: log errors
	s.Receive(ctx, "a", strings.NewReader("av"))
	s.Receive(ctx, "b", strings.NewReader("bv"))
	s.Receive(ctx, "d", strings.NewReader("dv"))
	s.Receive(ctx, "c", strings.NewReader("cv"))
	s.Receive(ctx, "e", strings.NewReader("ev"))
	s.Receive(ctx, "f", strings.NewReader("fv"))
	s.Receive(ctx, "g", strings.NewReader("gv"))
	s.Receive(ctx, "h", strings.NewReader("hv"))

	s.Remove(ctx, "b")
	t.Log(blobStr(s.Fetch(ctx, "b")))

	// TODO: test stat
	testEnumerate(ctx, t, s)
}

func testEnumerate(ctx context.Context, t *testing.T, s blob.Storage) {
	if err := s.Enumerate(context.Background(), "", func(ref string, size int64) error {
		header, err := blob.Peek(s.Fetch(ctx, ref), 1<<10)
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
