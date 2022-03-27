package localfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/koud-fi/pkg/blob"
)

const partFileExt = ".part"

type fileBlob struct {
	path string
}

func New(path string) blob.Blob {
	return &fileBlob{path: path}
}

func (b fileBlob) Open() (io.ReadCloser, error) {
	absPath, err := filepath.Abs(b.path)
	if err != nil {
		return nil, err
	}
	return os.Open(absPath)
}

func Write(path string, b blob.Blob, perm os.FileMode) error {
	rc, err := b.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	return WriteReader(path, rc, perm)
}

func WriteReader(path string, r io.Reader, perm os.FileMode) error {
	partPath := path + "." + strconv.FormatInt(time.Now().UnixNano(), 36) + partFileExt
	partFile, err := os.OpenFile(partPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer func() {
		if err := os.Remove(partPath); err != nil && !os.IsNotExist(err) {
			panic(fmt.Sprintf("failed to remove .part file: %v", err))
		}
	}()
	_, err = io.Copy(partFile, r)
	partFile.Close() // make sure that part file is always closed
	if err != nil {
		return err
	}
	return os.Rename(partPath, path)
}
