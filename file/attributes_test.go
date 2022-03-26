package file_test

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/file"
)

func TestAttributes(t *testing.T) {
	fsys := os.DirFS("../testdata")
	if err := fs.WalkDir(fsys, ".", func(path string, _ fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		attrs, err := file.ResolveAttrs(blob.FromFS(fsys, path),
			file.MediaAttrs())
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		t.Logf("%s %v", path, *attrs)
		return nil

	}); err != nil {
		t.Fatal(err)
	}
}
