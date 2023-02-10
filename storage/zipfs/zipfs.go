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

func (f *file) Stat() (fs.FileInfo, error) { return f, nil }

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

func (f *file) Name() string { return f.name }

func (f *file) Size() int64 {
	if f.file != nil {
		return int64(f.file.UncompressedSize)
	}
	return 0
}
func (f *file) Mode() fs.FileMode {
	if f.file != nil {
		return 0444
	}
	return os.ModeDir | 0555
}

func (f *file) ModTime() time.Time {
	if f.file != nil {
		return f.file.ModTime()
	}
	return time.Time{}
}

func (f *file) IsDir() bool { return f.file == nil }
func (f *file) Sys() any    { return nil }
