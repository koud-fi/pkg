package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/pk"
)

var registry = make(map[pk.Scheme]reflect.Type)

// Register calls are not thread-safe and should only be done on startup.
func Register(s pk.Scheme, v interface{}) {
	typ := reflect.TypeOf(v)
	switch typ.Kind() {
	case reflect.Struct, reflect.Slice:
	default:
		panic("schema.Register: v type must be struct or slice")
	}
	registry[s] = typ
}

func Decode(ref pk.Ref, b blob.Blob) (interface{}, []byte, error) {
	typ, ok := registry[ref.Scheme()]
	if !ok {
		return nil, nil, fmt.Errorf("%w: unknown schema; %s", os.ErrInvalid, ref.Scheme())
	}
	data, err := blob.Bytes(b)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading blob %s: %w", ref, err)
	}
	v := reflect.New(typ).Interface()
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, data, fmt.Errorf("error decoding schema: %w", err)
	}
	return v, data, nil
}
