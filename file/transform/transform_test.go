package transform_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/file/transform"
)

func Test(t *testing.T) {
	os.RemoveAll("temp")
	os.Mkdir("temp", os.FileMode(0700))

	fsys := os.DirFS("../../testdata")
	if err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		var (
			params, _ = transform.ParseParams("300x")
			outPath   = filepath.Join("temp", fmt.Sprintf("%s.%s.jpg", d.Name(), params))
			in        = blob.FromFS(fsys, path)
			out       = transform.ToImage(in, "", params)
		)
		if err := blob.WriteFile(outPath, out, os.FileMode(0600)); err != nil {
			t.Log("ERROR:", err)
		}
		return nil

	}); err != nil {
		t.Fatal(err)
	}
}
