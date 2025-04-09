package blob

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func FromFS(fsys fs.FS, name string) Reader {
	return Func(func() (io.ReadCloser, error) {
		return fsys.Open(name)
	})
}

func FromFile(name string) Reader {
	return Func(func() (io.ReadCloser, error) {
		absPath, err := filepath.Abs(name)
		if err != nil {
			return nil, err
		}
		return os.Open(absPath)
	})
}

// Save saves the blob to a file with default permissions, creating the directory if it doesn't exist.
func Save(path string, r Reader) error {
	const (
		filePerm = os.FileMode(0600)
		dirPerm  = os.FileMode(0700)
	)
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return fmt.Errorf("mkdirall: %w", err)
	}
	return WriteToFile(path, r, filePerm)
}

// WriteToFile writes the blob to a file with the given permissions.
func WriteToFile(path string, r Reader, perm os.FileMode) error {
	rc, err := r.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	return WriteIOReaderToFile(path, rc, perm)
}

// WriteIOReaderToFile writes the io.Reader to a file with the given permissions.
func WriteIOReaderToFile(path string, r io.Reader, perm os.FileMode) error {
	const partFileExt = ".part"
	partPath := path + "." + strconv.FormatInt(time.Now().UnixNano(), 36) + partFileExt
	partFile, err := os.OpenFile(partPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("openfile: %w", err)
	}
	defer func() {
		if err := os.Remove(partPath); err != nil && !os.IsNotExist(err) {
			panic(fmt.Sprintf("failed to remove .part file: %v", err))
		}
	}()
	_, err = io.Copy(partFile, r)
	partFile.Close() // Make sure that the part file is always closed.
	if err != nil {
		return fmt.Errorf("copy: %w", err)
	}
	return os.Rename(partPath, path)
}
