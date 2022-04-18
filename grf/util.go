package grf

import "encoding/json"

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
