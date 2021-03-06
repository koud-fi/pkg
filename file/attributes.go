package file

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/koud-fi/pkg/blob"
)

const defaultHeaderPeekSize = 512

type Attributes struct {
	Size        int64             `json:"size,omitempty"`
	ContentType string            `json:"contentType,omitempty"`
	Digest      map[string]string `json:"digest,omitempty"`
	MediaAttributes
	Info
}

type Info struct {
	ModTime *time.Time  `json:"modTime,omitempty"`
	Mode    os.FileMode `json:"-"`
	IsDir   bool        `json:"isDir,omitempty"`
}

func (a Attributes) Type() string {
	return strings.SplitN(a.ContentType, "/", 2)[0]
}

func (a Attributes) Ext() string {
	if exts, _ := mime.ExtensionsByType(a.ContentType); len(exts) > 0 {
		return exts[len(exts)-1]
	}
	return ""
}

type Option func(a *Attributes, b blob.Blob, contentType string) error

func ResolveAttrs(b blob.Blob, opts ...Option) (*Attributes, error) {
	var a Attributes

	switch b := b.(type) {
	case blob.BytesReader:
		buf := b.Bytes()
		a.Size = int64(len(buf))
		a.ContentType = http.DetectContentType(buf)
	default:
		rc, err := b.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		switch rc := rc.(type) {
		case fs.File:
			info, err := rc.Stat()
			if err != nil {
				return nil, err
			}
			a.Size = info.Size()

			if modTime := info.ModTime(); !modTime.IsZero() {
				a.ModTime = &modTime
			}
			a.Mode = info.Mode()
			a.IsDir = info.IsDir()

			if !a.IsDir {
				br := bufio.NewReaderSize(rc, defaultHeaderPeekSize)
				if a.ContentType, err = resolveContentType(br, info); err != nil {
					return nil, fmt.Errorf("failed to resolve content-type: %w", err)
				}
			}

		// case io.ReadSeekCloser:

		// TODO

		default:
			buf, err := io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("failed to read all data: %w", err)
			}
			a.Size = int64(len(buf))
			a.ContentType = http.DetectContentType(buf)

			b = blob.FromBytes(buf)
		}
	}
	for _, opt := range opts {
		if err := opt(&a, b, a.ContentType); err != nil {
			return nil, err
		}
	}
	return &a, nil
}

func resolveContentType(br *bufio.Reader, info os.FileInfo) (string, error) {
	if info != nil {
		if m := mime.TypeByExtension(filepath.Ext(info.Name())); m != "" {
			return m, nil
		}
	}
	b, err := br.Peek(defaultHeaderPeekSize)
	if err != nil && err != io.EOF {
		return "", err
	}
	if m := http.DetectContentType(b); m != "application/octet-stream" {
		return m, nil
	}
	return "", nil
}
