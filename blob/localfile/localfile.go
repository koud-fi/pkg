package localfile

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/koud-fi/pkg/blob"
)

const partFileExt = ".part"

func Write(path string, r blob.Reader, perm os.FileMode) error {
	rc, err := r.Open()
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
