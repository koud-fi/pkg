package proc

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/koud-fi/pkg/blob"
)

func WriteOutput(w io.Writer, out any) error {
	switch v := out.(type) {
	case nil:
		return nil
	case blob.Blob:
		return blob.WriteTo(w, v)
	}
	outType := reflect.TypeOf(out)
	switch outType.Kind() {
	case reflect.Chan:
		return writeChan(w, out)
	case reflect.Invalid, reflect.Func, reflect.UnsafePointer:
		panic(fmt.Sprintf("invalid output kind: %v", outType.Kind()))
	default:
		return writeValue(w, out)
	}
}

func writeChan(w io.Writer, c any) error {
	cVal := reflect.ValueOf(c)
	for {
		next, ok := cVal.Recv()
		if !ok {
			break
		}
		if err := writeValue(w, next.Interface()); err != nil {
			return err
		}
	}
	return nil
}

func writeValue(w io.Writer, v any) (err error) {
	switch v := v.(type) {
	case []byte:
		_, err = w.Write(v)
	case string, fmt.Stringer:
		_, err = fmt.Fprintln(w, v)
	default:
		if r, ok := v.(io.Reader); ok {
			if c, ok := r.(io.Closer); ok {
				defer c.Close()
			}
			_, err = io.Copy(w, r)
		} else {
			e := json.NewEncoder(w)
			e.SetIndent("", "\t")
			err = e.Encode(v)
		}
	}
	return
}
