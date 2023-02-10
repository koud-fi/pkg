package zipfs

import (
	"archive/zip"
	"io/fs"
	"sort"
)

type zipfs struct {
	*zip.ReadCloser
	list zipList
}

func New(rc *zip.ReadCloser) fs.FS {
	list := make(zipList, len(rc.File))
	copy(list, rc.File)
	sort.Sort(list)
	return &zipfs{rc, list}
}

func (fsys *zipfs) Open(name string) (fs.File, error) {

	// ???

	panic("TODO")
}

func (fsys *zipfs) ReadDir(n int) ([]fs.DirEntry, error) {

	// ???

	panic("TODO")
}

func (fsys *zipfs) Stat(name string) (fs.FileInfo, error) {

	// ???

	panic("TODO")
}

type zipList []*zip.File

func (z zipList) Len() int           { return len(z) }
func (z zipList) Less(i, j int) bool { return z[i].Name < z[j].Name }
func (z zipList) Swap(i, j int)      { z[i], z[j] = z[j], z[i] }
