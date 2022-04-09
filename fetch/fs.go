package fetch

import (
	"io/fs"
	"net/http"
)

type fetchfs struct {
	reqFn func(name string) *Request
}

func NewFS(reqFn func(name string) *Request) fs.StatFS {
	return &fetchfs{reqFn: reqFn}
}

func (fsys *fetchfs) Open(name string) (fs.File, error) {
	return fsys.reqFn(name).Method(http.MethodGet).OpenFile()
}

func (fsys *fetchfs) Stat(name string) (fs.FileInfo, error) {
	return fsys.reqFn(name).Method(http.MethodHead).Stat()
}
