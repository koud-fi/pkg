package grf

import "fmt"

type Node[T any] struct {
	ID      ID
	Type    NodeType
	Data    T
	Version int64
}

type NodeType string

func Lookup[T any](g *Graph, id ID) (*Node[T], error) {
	ti, s, err := g.parseID(id)
	if err != nil {
		return nil, err
	}
	ns, err := s.Node(ti.Type, id.localID())
	if err != nil {
		return nil, err
	}
	if len(ns) == 0 {
		return nil, fmt.Errorf("%w: %d", ErrNotFound, id)
	}
	v, err := unmarshal[T](ti.dataType, ns[0].Data)
	if err != nil {
		return nil, fmt.Errorf("invalid node data: %w", err)
	}
	return &Node[T]{
		ID:      id,
		Type:    ti.Type,
		Data:    v,
		Version: ns[0].Version,
	}, nil
}

func Add[T any](g *Graph, nt NodeType, v T) (*Node[T], error) {
	typeID, ti, shardID, s, err := g.resolveAddArgs(nt)
	if err != nil {
		return nil, err
	}
	localID, ver, err := s.AddNode(nt, marshal(v))
	if err != nil {
		return nil, err
	}
	return &Node[T]{
		ID:      newID(shardID, typeID, localID),
		Type:    ti.Type,
		Data:    v,
		Version: ver,
	}, nil
}

func Update[T any](g *Graph, id ID, fn func(T) (T, error)) (*Node[T], error) {
	n, err := Lookup[T](g, id)
	if err != nil {
		return nil, err
	}
	return update(g, n, fn)
}

func update[T any](g *Graph, n *Node[T], fn func(T) (T, error)) (*Node[T], error) {
	var err error
	if n.Data, err = fn(n.Data); err != nil {
		return nil, err
	}
	return n, g.shardForID(n.ID.shardID()).
		UpdateNode(n.Type, n.ID.localID(), marshal(n.Data), n.Version)
}

func Delete(g *Graph, id ID) error {
	ti, s, err := g.parseID(id)
	if err != nil {
		return err
	}

	// TODO: batch by store/type

	return s.DeleteNode(ti.Type, id.localID())
}

func (n Node[T]) String() string {
	return fmt.Sprintf("%d(%s)(%v) %v", n.ID, n.Type, n.Version, n.Data)
}
