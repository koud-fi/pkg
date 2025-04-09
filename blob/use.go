package blob

import "io"

func Use(r Reader, fn func(io.Reader) error) error {
	rc, err := r.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	return fn(rc)
}

func WriteTo(w io.Writer, r Reader) error {
	return Use(r, func(r io.Reader) error {
		_, err := io.Copy(w, r)
		return err
	})
}
