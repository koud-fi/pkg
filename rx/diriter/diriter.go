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

type config struct {
	hideFunc func(fs.DirEntry) bool
}

type Option func(*config)

func HideFunc(fn func(fs.DirEntry) bool) Option { return func(c *config) { c.hideFunc = fn } }

func DefaultHideFunc(d fs.DirEntry) bool {
	return strings.HasPrefix(d.Name(), ".")
}

func New(fsys fs.FS, root string, opt ...Option) rx.Iter[Entry] {
	c := config{
		hideFunc: DefaultHideFunc,
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
				} else if !c.hideFunc(topEntry) {
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
