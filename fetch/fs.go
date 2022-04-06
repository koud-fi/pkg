package fetch

import (
	"io/fs"
	"os"
	"strings"
)

var _ fs.StatFS = (*FS)(nil)

type FS struct {
	rootURL string
}

func NewFS(rootURL string) *FS {
	return &FS{rootURL: strings.TrimRight(rootURL, "/")}
}

type fetchfs struct {
	rootURL string
}

func (t *FS) Open(name string) (fs.File, error) {
	return Get(t.url(name)).OpenFile()
}

func (t *FS) Stat(name string) (os.FileInfo, error) {
	return Head(t.url(name)).Stat()
}

func (t *FS) url(name string) string {
	return strings.Join([]string{t.rootURL, name}, "/")
}
