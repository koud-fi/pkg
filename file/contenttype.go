package file

import (
	"bufio"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

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
