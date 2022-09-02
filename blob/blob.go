package blob

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"sync/atomic"
)

type Blob interface{ Open() (io.ReadCloser, error) }

type Func func() (io.ReadCloser, error)

func (fn Func) Open() (io.ReadCloser, error) { return fn() }

type ByteFunc func() ([]byte, error)

func (fn ByteFunc) Open() (io.ReadCloser, error) {
	buf, err := fn()
	if err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewReader(buf)), nil
}

func (fn ByteFunc) Bytes() ([]byte, error) { return fn() }

func FromBytes(buf []byte) Blob {
	return ByteFunc(func() ([]byte, error) { return buf, nil })
}

func FromString(s string) Blob { return FromBytes([]byte(s)) }

func FromReader(r io.Reader) Blob {
	var opened int32
	return Func(func() (io.ReadCloser, error) {
		if atomic.AddInt32(&opened, 1) > 1 {
			return nil, errors.New("multiple opens on io.Reader blob")
		}
		return io.NopCloser(r), nil
	})
}

func FromFS(fsys fs.FS, name string) Blob {
	return Func(func() (io.ReadCloser, error) { return fsys.Open(name) })
}

func Empty() Blob { return FromBytes(nil) }

type BytesReader interface {
	io.ReadCloser
	Bytes() []byte
}

func Bytes(b Blob) ([]byte, error) {
	rc, err := b.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	if br, ok := rc.(BytesReader); ok {
		return br.Bytes(), nil
	}
	return io.ReadAll(rc)
}

func Peek(b Blob, n int) ([]byte, error) {
	rc, err := b.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	if br, ok := rc.(BytesReader); ok {
		buf := br.Bytes()
		if len(buf) > n {
			return buf[:n], nil
		}
		return buf, nil
	}
	if br, ok := b.(interface{ Bytes() []byte }); ok {
		return br.Bytes(), nil
	}
	buf, err := bufio.NewReaderSize(rc, n).Peek(n)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return buf, nil
}

func String(b Blob) (string, error) {
	buf, err := Bytes(b)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func Error(b Blob) error {
	rc, err := b.Open()
	if err != nil {
		return err
	}
	return rc.Close()
}
