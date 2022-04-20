package grf

import (
	"fmt"
	"reflect"
	"time"
)

type Node[T any] struct {
	ID        ID
	Type      NodeType
	Data      T
	Timestamp time.Time
}

type NodeType string

func LookupAny(g *Graph, id ID) (*Node[any], error) {
	return lookup(g, id, unmarshalAny)
}

func Lookup[T any](g *Graph, id ID) (*Node[T], error) {
	return lookup(g, id, unmarshal[T])
}

func lookup[T any](
	g *Graph, id ID, unmarshal func(reflect.Type, []byte) (T, error),
) (*Node[T], error) {
	ti, s, err := g.parseID(id)
	if err != nil {
		return nil, fmt.Errorf("id parsing failed: %w", err)
	}
	ns, err := s.Node(ti.Type, id.localID())
	if err != nil {
		return nil, fmt.Errorf("store lookup failed: %w", err)
	}
	if len(ns) == 0 {
		return nil, fmt.Errorf("%w: %d", ErrNotFound, id)
	}
	nd := ns[0]
	data, err := unmarshal(ti.dataType, nd.Data)
	if err != nil {
		return nil, fmt.Errorf("data decoding failed: %w", err)
	}
	return &Node[T]{
		ID:        id,
		Type:      ti.Type,
		Data:      data,
		Timestamp: nd.Timestamp,
	}, nil
}

func (n Node[T]) String() string {
	ts := n.Timestamp.UTC().Format(time.RFC3339Nano)
	return fmt.Sprintf("%d(%s)(%v) %v", n.ID, n.Type, ts, n.Data)
}

/*
func (n *Node) Update(fn func(v any) (any, error)) *Node {
	if n.err != nil {
		return n
	}
	var v any
	if v, n.err = n.Data(); n.err != nil {
		return n
	}
	if v, n.err = fn(v); n.err != nil {
		return n
	}
	n.d.Data = marshal(v)
	n.err = n.s.UpdateNode(n.ti.Type, n.d.ID, n.d.Data)
	return n
}
*/
