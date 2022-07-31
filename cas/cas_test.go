package cas_test

import (
	"context"
	"crypto"
	"io/fs"
	"os"
	"testing"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/cas"
	"github.com/koud-fi/pkg/datastore"
	"github.com/koud-fi/pkg/file"
	"github.com/koud-fi/pkg/storage/localdisk"
)

func TestStorage(t *testing.T) {
	assert(t, os.RemoveAll("temp"))

	fsys := os.DirFS("../testdata")
	fileBlobs, err := localdisk.NewStorage("temp/file", localdisk.Buckets(1, 2))
	assert(t, err)
	metaBlobs, err := localdisk.NewStorage("temp/meta", localdisk.Buckets(1, 2))
	assert(t, err)

	var (
		ds = datastore.New[file.Attributes](metaBlobs, datastore.JSON())
		s  = cas.New(fileBlobs, ds, file.MediaAttrs(), file.Digests(crypto.MD5))
	)
	assert(t, fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		n, err := s.Add(context.Background(), blob.FromFS(fsys, path))
		if err != nil {
			return err
		}
		t.Logf("%s %v", path, n)
		return nil
	}))
}

func assert(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
