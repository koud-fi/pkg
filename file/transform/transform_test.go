package transform_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/file"
	"github.com/koud-fi/pkg/file/transform"
)

func TestToImage(t *testing.T) {
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
		in := blob.FromFS(fsys, path)
		attrs, err := file.ResolveAttrs(in, file.MediaAttrs())
		if err != nil {
			return err
		}
		for _, params := range transform.StdImagePreviewParamsList(attrs.MediaAttributes) {
			var (
				outPath = filepath.Join("temp", fmt.Sprintf("%s.%s.jpg", d.Name(), params))
				out     = transform.ToImage(in, params)
			)
			if err := blob.WriteFile(outPath, out, os.FileMode(0600)); err != nil {
				t.Log("ERROR:", err)
			}

			// TODO: actually check that the resulting images are of correct size

		}
		return nil

	}); err != nil {
		t.Fatal(err)
	}
}
