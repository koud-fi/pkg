package file

import (
	"bufio"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/koud-fi/pkg/file/format/raw"
)

func init() {
	mime.AddExtensionType(".raf", raw.RAFMime)
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
