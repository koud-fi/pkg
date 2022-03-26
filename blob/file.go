package blob

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const partFileExt = ".part"

type fileBlob struct {
	path string
}

func FromFile(path string) Blob {
	return &fileBlob{path: path}
}

func (b fileBlob) Open() (io.ReadCloser, error) {
	absPath, err := filepath.Abs(b.path)
	if err != nil {
		return nil, err
	}
	return os.Open(absPath)
}

func WriteFile(path string, b Blob, perm os.FileMode) error {
	rc, err := b.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

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
	_, err = io.Copy(partFile, rc)
	partFile.Close() // make sure that part file is always closed
	if err != nil {
		return err
	}
	return os.Rename(partPath, path)
}
