package merge_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/koud-fi/pkg/merge"
)

func TestAssing(t *testing.T) {
	type (
		InnerData struct {
			C bool
			D string
			S []int
		}
		Data struct {
			A      int
			B      string
			b      string
			Nested InnerData
		}
	)
	var (
		assFn = merge.NewAssignFunc(reflect.TypeOf(&Data{}))
		dst   Data
	)
	if err := assFn(reflect.ValueOf(&dst), Data{
		A: 42,
		B: "Hello, world?",
		b: "hello, world?",
		Nested: InnerData{
			C: true,
			D: "derp",
			S: []int{1, 2, 3},
		},
	}); err != nil {
		t.Fatal(err)
	}
	out, _ := json.MarshalIndent(dst, "", "\t")
	t.Log(string(out))
}
