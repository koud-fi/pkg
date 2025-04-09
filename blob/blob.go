package blob

import (
	"bufio"
	"errors"
	"io"
	"sync/atomic"
)

type Reader interface{ Open() (io.ReadCloser, error) }
type Blob = Reader // TODO: refactor everything to use Reader instead of Blob

type Func func() (io.ReadCloser, error)

func (fn Func) Open() (io.ReadCloser, error) { return fn() }

func FromString(s string) Reader { return FromBytes([]byte(s)) }

func FromReader(r io.Reader) Reader {
	var opened int32
	return Func(func() (io.ReadCloser, error) {
		if atomic.AddInt32(&opened, 1) > 1 {
			return nil, errors.New("multiple opens on io.Reader blob")
		}
		return io.NopCloser(r), nil
	})
}

func Empty() Reader { return FromBytes(nil) }

func ReadAt(p []byte, r Reader, n int64) (int, error) {
	rc, err := r.Open()
	if err != nil {
		return 0, err
	}
	defer rc.Close()

	switch r := rc.(type) {
	case io.ReaderAt:
		return r.ReadAt(p, n)

	case io.ReadSeeker:
		if _, err := r.Seek(n, io.SeekStart); err != nil {
			return 0, err
		}
		return r.Read(p)

	default:
		br := bufio.NewReader(r)
		if _, err := br.Discard(int(n)); err != nil {
			return 0, err
		}
		return br.Read(p)
	}
}

func Peek(r Reader, n int) ([]byte, error) {
	rc, err := r.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	if br, ok := rc.(interface {
		Bytes() []byte
	}); ok {
		buf := br.Bytes()
		if len(buf) > n {
			return buf[:n], nil
		}
		return buf, nil
	}
	if br, ok := r.(interface{ Bytes() []byte }); ok {
		return br.Bytes(), nil
	}
	buf, err := bufio.NewReaderSize(rc, n).Peek(n)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return buf, nil
}

func String(r Reader) (string, error) {
	buf, err := Bytes(r)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func Error(r Reader) error {
	rc, err := r.Open()
	if err != nil {
		return err
	}
	return rc.Close()
}
