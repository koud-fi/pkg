package diriter

import (
	"io/fs"
	"path"
	"strings"

	"github.com/koud-fi/pkg/rx"
)

const defaultBatchSize = 1 << 6

type Entry struct {
	path string
	fs.DirEntry
}

func (e Entry) Path() string { return e.path }

func New(fsys fs.FS, root string) rx.Iter[Entry] {
	type dirInfo struct {
		path    string
		entries []fs.DirEntry
	}
	var (
		dirs []dirInfo
		init bool
	)
	return rx.FuncIter(func() ([]Entry, bool, error) {
		if !init {
			dir, err := fs.ReadDir(fsys, path.Clean(root))
			if err != nil {
				return nil, false, err
			}
			dirs = []dirInfo{{path: root, entries: dir}}
			init = true
		}
		out := make([]Entry, 0, defaultBatchSize)
		for len(out) <= defaultBatchSize {
			if len(dirs) == 0 {
				break
			}
			i := len(dirs) - 1
			if len(dirs[i].entries) == 0 {
				dirs = dirs[:len(dirs)-1]
			} else {
				var (
					topEntry = dirs[i].entries[len(dirs[i].entries)-1]
					topPath  = path.Join(dirs[i].path, topEntry.Name())
				)
				if topEntry.IsDir() {
					dir, err := fs.ReadDir(fsys, topPath)
					if err != nil {
						return nil, false, err
					}
					dirs = append(dirs, dirInfo{
						path:    topPath,
						entries: dir,
					})
				} else {
					out = append(out, Entry{
						path:     topPath,
						DirEntry: topEntry,
					})
				}
				dirs[i].entries = dirs[i].entries[:len(dirs[i].entries)-1]
			}
		}
		return out, len(dirs) > 0, nil
	})
}

func IsHidden(e Entry) bool {
	return strings.HasPrefix(e.Name(), ".") // TODO: check for other common hidden file patters
}
