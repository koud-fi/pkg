package blob

import "encoding/json"

type MarshalFunc func(v any) ([]byte, error)
type UnmarshalFunc func(data []byte, v any) error

func Marshal(fn MarshalFunc, v any) Reader {
	return ByteFunc(func() ([]byte, error) { return fn(v) })
}

func MarshalJSON(v any) Reader {
	return Marshal(json.Marshal, v)
}

func MarshalJSONIndent(v any, prefix, indent string) Reader {
	return ByteFunc(func() ([]byte, error) {
		return json.MarshalIndent(v, prefix, indent)
	})
}

func MarshalJSONPretty(v any) Reader {
	return MarshalJSONIndent(v, "", "\t")
}

func Unmarshal(fn UnmarshalFunc, r Reader, v any) error {
	buf, err := Bytes(r)
	if err != nil {
		return err
	}
	return fn(buf, v)
}

func UnmarshalJSON(r Reader, v any) error {
	return Unmarshal(json.Unmarshal, r, v)
}

func UnmarshalJSONValue[T any](r Reader) (T, error) {
	var v T
	return v, UnmarshalJSON(r, &v)
}
