package fetch

import (
	"io"
	"io/fs"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

type file struct {
	*fileInfo
	body io.ReadCloser
}

func (f *file) Stat() (fs.FileInfo, error) { return f.fileInfo, nil }
func (f *file) Read(p []byte) (int, error) { return f.body.Read(p) }
func (f *file) Close() error               { return f.body.Close() }

func (f *file) ReadDir(n int) ([]fs.DirEntry, error) {
	if f.dr != nil {
		return f.dr.ReadDir(f, f.header, n)
	}
	return nil, fs.ErrInvalid
}

type fileInfo struct {
	url    *url.URL
	header http.Header
	dr     DirReader
}

func (fi fileInfo) Name() string {
	name := path.Base(fi.url.Path)
	if name == "" || name == "/" {
		name = "."
	}
	if path.Ext(name) == "" {
		exts, _ := mime.ExtensionsByType(fi.header.Get("Content-Type"))
		if len(exts) > 0 {
			name += exts[0]
		}
	}
	return name
}

func (fi fileInfo) Size() int64 {
	size, _ := strconv.ParseInt(fi.header.Get("Content-Length"), 10, 64)
	return size
}

func (fi fileInfo) Mode() fs.FileMode { return fs.FileMode(0700) }

func (fi fileInfo) ModTime() time.Time {
	modTime, _ := time.Parse(http.TimeFormat, fi.header.Get("Last-Modified"))
	return modTime
}

func (fi fileInfo) IsDir() bool {
	if fi.dr != nil {
		return fi.dr.IsDir(fi.header)
	}
	return false
}

func (t fileInfo) Sys() any { return nil }
