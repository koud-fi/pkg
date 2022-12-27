package rq_test

import (
	"encoding/json"
	"testing"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/blob/localfile"
	"github.com/koud-fi/pkg/rq"
)

func TestMap(t *testing.T) {
	var data map[string]any
	if err := blob.Unmarshal(json.Unmarshal, localfile.New("testdata.json"), &data); err != nil {
		t.Fatal(err)
	}
	n := rq.MapNode[any](data)
	t.Log(rq.Attr(n, "a"))
	t.Log(rq.Attr(n, "b"))
}
