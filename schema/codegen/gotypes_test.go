package codegen_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/koud-fi/pkg/schema"
	"github.com/koud-fi/pkg/schema/codegen"
)

func TestGoTypes(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	if err := codegen.GoTypes(buf, schema.Schema{
		Definitions: map[string]schema.Type{
			"Thing": schema.ResolveType[struct {
				A string
				B int
				C float64
				D bool `json:"_d"`
			}](),
		},
	}); err != nil {
		log.Fatal(err)
	}
	t.Log(buf.String())
}
