package blob

import (
	"bytes"
	"io"
)

func Bytes(b Reader) ([]byte, error) {
	rc, err := b.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	if br, ok := rc.(interface {
		io.ReadCloser
		Bytes() []byte
	}); ok {
		return br.Bytes(), nil
	}
	return io.ReadAll(rc)
}

type ByteFunc func() ([]byte, error)

func (fn ByteFunc) Open() (io.ReadCloser, error) {
	buf, err := fn()
	if err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewReader(buf)), nil
}

func (fn ByteFunc) Bytes() ([]byte, error) {
	return fn()
}

func FromBytes(buf []byte) Reader {
	return ByteFunc(func() ([]byte, error) { return buf, nil })
}
