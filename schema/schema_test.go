package schema_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/schema"
)

//go:embed test.json
var testJSON []byte

func TestResolveFromType(t *testing.T) {
	t.Log(schema.Resolve[int]())
	t.Log(schema.Resolve[struct {
		Value   string    `json:"value"`
		Numbers []float64 `json:"nums"`
		Things  []struct {
			ID   int64  `json:"id"`
			Name string `json:"name,omitempty"`
		} `json:"things"`
	}]().ExampleJSON())

}

func TestResolveFromValue(t *testing.T) {
	var v any
	if err := blob.Unmarshal(json.Unmarshal, blob.FromBytes(testJSON), &v); err != nil {
		t.Fatal(err)
	}
	t.Log(schema.ResolveFromValue(v).ExampleJSON())
}
