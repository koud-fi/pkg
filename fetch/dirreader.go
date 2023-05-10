package fetch

import (
	"io/fs"
	"net/http"
)

type DirReader interface {
	IsDir(h http.Header) bool
	ReadDir(f fs.File, h http.Header, n int) ([]fs.DirEntry, error)
}
