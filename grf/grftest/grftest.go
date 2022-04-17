package grftest

import (
	"testing"
	"time"

	"github.com/koud-fi/pkg/grf"
)

var types = []grf.NodeType{"type1", "type2"}

func Test(t *testing.T, s ...grf.Store) {
	g := grf.New(nil, s...)
	g.Register(types[0], 1)
	g.Register(types[1], 2)

	ns := []*grf.Node{
		assert1[*grf.Node](t)(g.AddNode(types[0], "Hello,")),
		assert1[*grf.Node](t)(g.AddNode(types[0], "World!")),
		assert1[*grf.Node](t)(g.AddNode(types[1], 42)),
	}
	assert(t, g.DeleteNode(ns[0].ID()))
	assert(t, g.UpdateNode(ns[1].ID(), "World?"))

	for _, n := range ns {
		t.Log(g.Node(n.ID()))
	}
	t.Log(g.Node(grf.ID(-1)))

	assert(t, g.SetEdge(grf.NewEdge(ns[1].ID(), "type1", ns[2].ID(), 1, nil)))
	assert(t, g.SetEdge(grf.NewEdge(ns[2].ID(), "type1", ns[1].ID(), 1, nil)))
	assert(t, g.SetEdge(grf.NewEdge(ns[1].ID(), "type2", ns[2].ID(), 1, nil)))

	for i, s := range s {
		t.Logf("shard %d:", i+1)
		for _, nt := range types {
			t.Logf("\t%s:", nt)
			for _, n := range assert1[[]grf.NodeData](t)(s.NodeRange(nt, 0, 10)) {
				t.Logf("\t\t%d %s %s", n.ID, n.Timestamp.Format(time.RFC3339Nano), string(n.Data))
			}
		}
	}
}

func assert(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func assert1[T any](t *testing.T) func(T, error) T {
	return func(v T, err error) T {
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
}
