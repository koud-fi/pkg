package grf

import "encoding/json"

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
