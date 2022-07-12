package blob

import (
	"bufio"
	"io"

	"github.com/koud-fi/pkg/rx"
)

func LineIter(b Blob) rx.Iter[string] {
	var (
		rc io.ReadCloser
		s  *bufio.Scanner
	)
	return rx.WithClose(rx.FuncIter(func() ([]string, bool, error) {
		if s == nil {
			var err error
			if rc, err = b.Open(); err != nil {
				return nil, false, err
			}
			s = bufio.NewScanner(rc)
		}
		if !s.Scan() {
			return nil, false, nil
		}
		return []string{s.Text()}, true, nil

	}), func() error {
		defer rc.Close()
		return s.Err()
	})
}

func ReadLines(b Blob, fn func(string) error) error {
	return rx.ForEach(LineIter(b), fn)
}
