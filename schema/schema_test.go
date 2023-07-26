package schema_test

import (
	"testing"

	"github.com/koud-fi/pkg/schema"
)

func TestResolve(t *testing.T) {
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
