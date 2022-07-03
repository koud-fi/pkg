package diriter

import (
	"io/fs"
	"path"

	"github.com/koud-fi/pkg/rx"
)

const defaultBatchSize = 1 << 6

type Entry struct {
	Path string
	fs.DirEntry
}

type dirInfo struct {
	path    string
	entries []fs.DirEntry
}

func New(fsys fs.FS, root string) rx.Iter[Entry] {
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
			topDir := &dirs[len(dirs)-1]
			if len(topDir.entries) == 0 {
				dirs = dirs[:len(dirs)-1]
			} else {
				var (
					topEntry = topDir.entries[len(topDir.entries)-1]
					topPath  = path.Join(topDir.path, topEntry.Name())
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
						Path:     topPath,
						DirEntry: topEntry,
					})
				}
				topDir.entries = topDir.entries[:len(topDir.entries)-1]
			}
		}
		return out, len(dirs) > 0, nil
	})
}

func Paths(it rx.Iter[Entry]) rx.Iter[string] {
	return rx.Map(it, func(e Entry) string { return e.Path })
}
