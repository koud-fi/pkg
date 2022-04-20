package grf

import (
	"encoding/json"
	"reflect"
)

func marshal(v any) []byte {
	if v == nil {
		return nil
	}
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

func unmarshal[T any](typ reflect.Type, data []byte) (v T, err error) {
	if data == nil {
		return
	}
	if reflect.TypeOf(v) == nil {
		if typ != nil {
			p := reflect.New(typ)
			if data != nil {
				err = json.Unmarshal(data, p.Interface())
			}
			v = p.Elem().Interface().(T)
		} else {
			err = json.Unmarshal(data, &v)
		}
	} else {
		err = json.Unmarshal(data, &v)
	}
	return
}
