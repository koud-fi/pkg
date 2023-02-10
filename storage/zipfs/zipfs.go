package zipfs

import (
	"archive/zip"
	"io/fs"
	"path"
	"sort"
	"strings"
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
	f := file{name: path.Clean(name)}
	if f.name == "." {
		return &f, nil
	}
	i, exact := fsys.list.lookup(f.name)
	if i < 0 {
		return nil, fs.ErrNotExist // TODO: return os.PathError
	}
	if exact {
		f.file = fsys.list[i]
	}
	return &f, nil
}

// TODO: ReadDirFS
// TODO: StatFS

type zipList []*zip.File

func (z zipList) Len() int           { return len(z) }
func (z zipList) Less(i, j int) bool { return z[i].Name < z[j].Name }
func (z zipList) Swap(i, j int)      { z[i], z[j] = z[j], z[i] }

func (z zipList) lookup(name string) (index int, exact bool) {
	i := sort.Search(len(z), func(i int) bool {
		return name <= z[i].Name
	})
	if i >= len(z) {
		return -1, false
	}
	if z[i].Name == name {
		return i, true
	}
	z = z[i:]
	name += "/"
	j := sort.Search(len(z), func(i int) bool {
		return name <= z[i].Name
	})
	if j >= len(z) {
		return -1, false
	}
	if strings.HasPrefix(z[j].Name, name) {
		return i + j, false
	}
	return -1, false
}
