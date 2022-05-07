package proc

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/koud-fi/pkg/blob"
)

func OutputBlob(out any) blob.Blob {
	switch v := out.(type) {
	case nil:
		return blob.Empty()
	case blob.Blob:
		return v
	}
	outType := reflect.TypeOf(out)
	switch outType.Kind() {

	// TODO: reflect.Chan

	case reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Invalid:
		panic(fmt.Sprintf("invalid output kind: %v", outType.Kind()))
	default:
		return valueBlob(out)
	}
}

func valueBlob(v any) blob.Blob {
	switch v := v.(type) {
	case []byte:
		return blob.FromBytes(v)
	case string:
		return blob.FromString(v)
	case fmt.Stringer:
		return blob.FromString(v.String())
	default:
		return blob.Marshal(json.Marshal, v)
	}
}
