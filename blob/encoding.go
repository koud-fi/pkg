package blob

type MarshalFunc func(v interface{}) ([]byte, error)
type UnmarshalFunc func(data []byte, v interface{}) error

func Marshal(fn MarshalFunc, v interface{}) Blob {
	return ByteFunc(func() ([]byte, error) { return fn(v) })
}

func Unmarshal(fn UnmarshalFunc, b Blob, v interface{}) error {
	buf, err := Bytes(b)
	if err != nil {
		return err
	}
	return fn(buf, v)
}
