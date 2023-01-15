package serializable

import (
	"bytes"
	"io"
)

type Value[T interface {
	Encode(io.Writer) error
	Decode(io.Reader) error
}] struct{ Data T }

func (v Value[T]) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil) // TODO: size hint for prealloc?
	buf.WriteByte('"')
	if err := v.Data.Encode(buf); err != nil {
		return nil, err
	}
	buf.WriteByte('"')
	return buf.Bytes(), nil
}

func (v Value[T]) UnmarshalJSON(data []byte) error {
	n := len(data)
	if n >= 2 && data[0] == '"' && data[n-1] == '"' {
		data = data[1 : n-1]
	}
	return v.Data.Decode(bytes.NewReader(data))
}
