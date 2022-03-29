package blob

type MarshalFunc func(v any) ([]byte, error)
type UnmarshalFunc func(data []byte, v any) error

func Marshal(fn MarshalFunc, v any) Blob {
	return ByteFunc(func() ([]byte, error) { return fn(v) })
}

func Unmarshal(fn UnmarshalFunc, b Blob, v any) error {
	buf, err := Bytes(b)
	if err != nil {
		return err
	}
	return fn(buf, v)
}
