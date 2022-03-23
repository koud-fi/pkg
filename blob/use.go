package blob

import (
	"bufio"
	"io"
)

func Use(b Blob, fn func(io.Reader) error) error {
	rc, err := b.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	return fn(rc)
}

func ReadLines(b Blob, fn func(string) error) error {
	return Use(b, func(r io.Reader) error {
		sr := bufio.NewScanner(r)
		for sr.Scan() {
			if err := fn(sr.Text()); err != nil {
				return err
			}
		}
		return sr.Err()
	})
}

func WriteTo(w io.Writer, b Blob) error {
	return Use(b, func(r io.Reader) error {
		_, err := io.Copy(w, r)
		return err
	})
}
