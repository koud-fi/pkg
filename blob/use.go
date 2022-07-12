package blob

import "io"

func Use(b Blob, fn func(io.Reader) error) error {
	rc, err := b.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	return fn(rc)
}

func WriteTo(w io.Writer, b Blob) error {
	return Use(b, func(r io.Reader) error {
		_, err := io.Copy(w, r)
		return err
	})
}
