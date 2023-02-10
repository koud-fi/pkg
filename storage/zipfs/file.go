package zipfs

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"strings"
	"time"
)

type file struct {
	name string
	file *zip.File
	dir  []*zip.File
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
	var (
		dir     []fs.DirEntry
		dirName = f.name
	)
	if dirName == "." {
		dirName = ""
	} else {
		dirName += "/"
	}
	var prevName string
	for _, e := range f.dir {
		if !strings.HasPrefix(e.Name, dirName) {
			break
		}
		f := file{
			name: e.Name[len(dirName):],
			file: e,
		}
		if i := strings.IndexRune(f.name, '/'); i >= 0 {
			f.name = f.name[0:i]
			f.file = nil
		}
		if f.name != prevName {
			dir = append(dir, (*fileInfo)(&f))
			prevName = f.name
		}
	}
	return dir, nil
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

func (fi *fileInfo) Type() fs.FileMode {
	if fi.file != nil {
		return fi.file.Mode().Type()
	}
	return os.ModeDir
}

func (fi *fileInfo) Info() (fs.FileInfo, error) { return fi, nil }
