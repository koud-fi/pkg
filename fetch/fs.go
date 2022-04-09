package fetch

import (
	"io/fs"
	"net/http"
)

var _ fs.StatFS = (*FS)(nil)

type FS struct {
	reqFn func(name string) *Request
}

func NewFS(reqFn func(name string) *Request) *FS {
	return &FS{reqFn: reqFn}
}

func (fsys *FS) Open(name string) (fs.File, error) {
	return fsys.reqFn(name).Method(http.MethodGet).OpenFile()
}

func (fsys *FS) Stat(name string) (fs.FileInfo, error) {
	return fsys.reqFn(name).Method(http.MethodHead).Stat()
}
