package rq_test

import (
	"testing"

	"github.com/koud-fi/pkg/rq"
)

func TestMap(t *testing.T) {
	n := rq.MapNode[any](map[string]any{"a": 42})
	t.Log(rq.Attr(n, "a"))
}
