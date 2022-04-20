package grftest

import (
	"testing"
	"time"

	"github.com/koud-fi/pkg/grf"
)

var types = []grf.NodeType{"type1", "type2"}

func Test(t *testing.T, s ...grf.Store) {
	g := grf.New(nil, s...)
	g.Register(
		grf.TypeInfo{
			Type:     types[0],
			DataType: "",
			Edges: []grf.EdgeTypeInfo{
				{Type: "type1"},
				{Type: "type2"},
			},
		},
		grf.TypeInfo{
			Type:     types[1],
			DataType: []int{},
			Edges: []grf.EdgeTypeInfo{
				{Type: "type1"},
			},
		})

	ns := []*grf.Node[any]{
		assert1[*grf.Node[any]](t)(grf.Add[any](g, types[0], "Hello,")),
		assert1[*grf.Node[any]](t)(grf.Add[any](g, types[0], "World!")),
		assert1[*grf.Node[any]](t)(grf.Add[any](g, types[0], "A")),
		assert1[*grf.Node[any]](t)(grf.Add[any](g, types[0], "B")),
		assert1[*grf.Node[any]](t)(grf.Add[any](g, types[0], "C")),
		assert1[*grf.Node[any]](t)(grf.Add[any](g, types[1], []int{42, 69})),
	}
	//assert(t, g.DeleteNode(ns[0].ID()))
	//assert(t, g.Node(ns[1].ID()).Update(func(_ any) (any, error) { return "World?", nil }).Err())
	for _, n := range ns {
		t.Log(grf.Lookup[any](g, n.ID))
		t.Log(grf.Lookup[string](g, n.ID))
		t.Log(grf.Lookup[[]int](g, n.ID))
	}
	t.Log(grf.Lookup[any](g, -1))

	/*
		assert(t, g.SetEdge(grf.NewEdge(ns[1].ID(), "type1", ns[2].ID(), 1, nil)))
		assert(t, g.SetEdge(grf.NewEdge(ns[2].ID(), "type1", ns[1].ID(), 1, nil)))
		assert(t, g.SetEdge(grf.NewEdge(ns[5].ID(), "type1", ns[1].ID(), 1, nil)))
		assert(t, g.SetEdge(grf.NewEdge(ns[1].ID(), "type1", ns[2].ID(), 2, nil)))
		assert(t, g.SetEdge(grf.NewEdge(ns[1].ID(), "type1", ns[3].ID(), 3, nil)))
		assert(t, g.SetEdge(grf.NewEdge(ns[1].ID(), "type1", ns[4].ID(), 4, nil)))
		assert(t, g.SetEdge(grf.NewEdge(ns[1].ID(), "type1", ns[5].ID(), 5, nil)))
		assert(t, g.SetEdge(grf.NewEdge(ns[1].ID(), "type2", ns[5].ID(), 1, nil)))
	*/

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
