package grf

import (
	"encoding/json"
	"reflect"
)

func marshal(v any) []byte {
	if v == nil {
		return nil
	}

	// TODO: have typeInfo.dataType passed in and validate v against it

	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

func unmarshal[T any](_ reflect.Type, data []byte) (v T, err error) {
	if data != nil {
		err = json.Unmarshal(data, &v)
	}
	return v, err
}

func unmarshalAny(typ reflect.Type, data []byte) (v any, err error) {
	if typ != nil {
		p := reflect.New(typ)
		if data != nil {
			err = json.Unmarshal(data, p.Interface())
		}
		v = p.Elem().Interface()
	} else if data != nil {
		err = json.Unmarshal(data, &v)
	}
	return
}
