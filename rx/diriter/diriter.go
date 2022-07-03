package diriter

import (
	"io/fs"
	"path"
	"strings"

	"github.com/koud-fi/pkg/rx"
)

const defaultBatchSize = 1 << 6

type Entry struct {
	Path string
	fs.DirEntry
}

type config struct {
	hideFunc func(string) bool
}

type Option func(*config)

func HideFunc(fn func(name string) bool) Option { return func(c *config) { c.hideFunc = fn } }

func New(fsys fs.FS, root string, opt ...Option) rx.Iter[Entry] {
	c := config{
		hideFunc: func(name string) bool { return strings.HasPrefix(name, ".") },
	}
	for _, opt := range opt {
		opt(&c)
	}
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
				} else if !c.hideFunc(topEntry.Name()) {
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
