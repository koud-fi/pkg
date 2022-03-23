package blob_test

import (
	"encoding/json"
	"testing"

	"github.com/koud-fi/pkg/blob"
)

type dataType struct{ N int }

func TestMarshal(t *testing.T) {
	s, err := blob.String(blob.Marshal(json.Marshal, dataType{42}))
	if err != nil {
		t.Fatal(err)
	}
	if s != `{"N":42}` {
		t.Fatal("mismatch:", s)
	}
}

func TestUnmarshal(t *testing.T) {
	var data dataType
	if err := blob.Unmarshal(json.Unmarshal, blob.FromString(`{"N":42}`), &data); err != nil {
		t.Fatal(err)
	}
	if data.N != 42 {
		t.Fatal("mismatch")
	}
}
