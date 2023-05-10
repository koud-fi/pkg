package fetch

import (
	"bufio"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultDirContentType = "text/vnd.directory; charset=utf-8"

type DirReader interface {
	IsDir(h http.Header) bool
	ReadDir(f fs.File, h http.Header, n int) ([]fs.DirEntry, error)
}

type dirReader struct{}

func DefaultDirReader() DirReader { return &dirReader{} }

func (dr dirReader) IsDir(h http.Header) bool {
	return h.Get("Content-Type") == defaultDirContentType
}

func (dr dirReader) ReadDir(f fs.File, _ http.Header, n int) ([]fs.DirEntry, error) {
	var (
		r   = bufio.NewReader(f)
		dir []fs.DirEntry
	)
	for i := 0; n <= 0 || i < n; i++ {
		s, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		d, err := decodeDirEntry(s)
		if err != nil {
			return nil, err
		}
		dir = append(dir, d)
	}
	return dir, nil
}

type direntry []string

func decodeDirEntry(s string) (fs.DirEntry, error) {
	parts := strings.Fields(s)
	if parts[0] == "" {
		return nil, errors.New("invalid dir entry")
	}
	return direntry(parts), nil
}

func (d direntry) Name() string               { return strings.TrimSuffix(d.part(0), "/") }
func (d direntry) IsDir() bool                { return strings.HasSuffix(d.part(0), "/") }
func (d direntry) Type() fs.FileMode          { return d.Mode().Type() }
func (d direntry) Info() (fs.FileInfo, error) { return d, nil }

func (d direntry) Size() int64 {
	size, _ := strconv.ParseInt(d.part(1), 10, 64)
	return size
}

func (d direntry) Mode() os.FileMode { return os.FileMode(0700) }

func (d direntry) ModTime() time.Time {
	modTime, _ := time.Parse("2006-01-02T15:04:05", d.part(2))
	return modTime
}

func (d direntry) Sys() any { return nil }

func (d direntry) part(i int) string {
	if i < 0 || i >= len(d) {
		return ""
	}
	return d[i]
}
