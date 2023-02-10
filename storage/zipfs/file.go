package zipfs

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"time"
)

type file struct {
	name string
	file *zip.File
	rc   io.ReadCloser
}

func (f *file) Stat() (fs.FileInfo, error) { return (*fileInfo)(f), nil }

func (f *file) Read(p []byte) (_ int, err error) {
	if f.rc == nil {
		if f.rc, err = f.file.Open(); err != nil {
			return
		}
	}
	return f.rc.Read(p)
}

func (f *file) Close() error {
	if f.rc != nil {
		return f.rc.Close()
	}
	return nil
}

func (f *file) ReadDir(n int) ([]fs.DirEntry, error) {

	// ???

	panic("TODO")
}

type fileInfo file

func (fi *fileInfo) Name() string { return fi.name }

func (fi *fileInfo) Size() int64 {
	if fi.file != nil {
		return int64(fi.file.UncompressedSize)
	}
	return 0
}
func (fi *fileInfo) Mode() fs.FileMode {
	if fi.file != nil {
		return 0444
	}
	return os.ModeDir | 0555
}

func (fi *fileInfo) ModTime() time.Time {
	if fi.file != nil {
		return fi.file.ModTime()
	}
	return time.Time{}
}

func (fi *fileInfo) IsDir() bool { return fi.file == nil }
func (*fileInfo) Sys() any       { return nil }
