package datastore

import (
	"encoding/json"

	"github.com/koud-fi/pkg/blob"
)

type Codec struct {
	Marshal   blob.MarshalFunc
	Unmarshal blob.UnmarshalFunc
}

func JSON() Codec {
	return Codec{
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
	}
}
