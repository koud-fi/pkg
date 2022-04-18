package grf

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidType     = errors.New("invalid type")
	ErrInvalidEdgeType = errors.New("invalid edge type")
	ErrAlreadyExists   = errors.New("already exists")
)

type Graph struct {
	m       Mapper
	shards  []Store
	counter int64
	schema
}

func New(m Mapper, s ...Store) *Graph {
	if len(s) == 0 {
		panic("no stores")
	}
	return &Graph{
		m:      m,
		shards: s,
		schema: schema{typeMap: make(map[NodeType]typeID)},
	}
}

func (g *Graph) Register(ti ...TypeInfo) {
	for _, ti := range ti {
		g.register(ti)
	}
}

func (g *Graph) Node(id ID) (*Node, error) {
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
	return &Node{
		id: id,
		t:  ti.Type,
		d:  ns[0],
	}, nil
}

func (g *Graph) MappedNode(nt NodeType, key string) (*Node, error) {
	id, err := g.m.Map(nt, key)
	if err != nil {
		return nil, err
	}
	return g.Node(id)
}

func (g *Graph) AddNode(nt NodeType, v any) (*Node, error) {
	typeID := g.typeMap[nt]
	if _, ok := g.typeInfo(typeID); !ok {
		return nil, ErrInvalidType
	}
	var (
		data     = marshal(v)
		shardNum = int(atomic.AddInt64(&g.counter, 1) % int64(len(g.shards)))
		shard    = g.shards[shardNum]
		shardID  = shardID(shardNum + 1)
	)
	localID, ts, err := shard.AddNode(nt, data)
	if err != nil {
		return nil, err
	}
	return &Node{
		id: newID(shardID, typeID, localID),
		t:  nt,
		d: NodeData{
			ID:        localID,
			Data:      data,
			Timestamp: ts,
		},
	}, nil
}

func (g *Graph) AddMappedNode(nt NodeType, key string, v any) (*Node, error) {
	var n *Node
	if id, err := g.m.Map(nt, key); err == nil {
		if n, err = g.Node(id); err != nil {
			return nil, err
		}
		if n.d.Data != nil {
			return nil, ErrAlreadyExists
		}
	} else if err == ErrNotFound {
		if n, err = g.AddNode(nt, nil); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	if err := g.m.SetMapping(nt, key, n.id); err != nil {
		return nil, err
	}
	return n, g.UpdateNode(n.ID(), v)
}

func (g *Graph) UpdateNode(id ID, v any) error {
	ti, s, err := g.parseID(id)
	if err != nil {
		return err
	}
	return s.UpdateNode(ti.Type, id.localID(), marshal(v))
}

func (g *Graph) DeleteNode(id ID) error {
	ti, s, err := g.parseID(id)
	if err != nil {
		return err
	}

	// TODO: support batching

	return s.DeleteNode(ti.Type, id.localID())
}

func (g *Graph) SetEdge(e Edge) error {
	ti, s, err := g.parseID(e.from)
	if err != nil {
		return err
	}
	if etID, ok := ti.edgeTypeMap[e.t]; ok {
		e.d.TypeID = etID
	} else {
		return ErrInvalidEdgeType
	}

	// TODO: support batching

	return s.SetEdge(ti.Type, e.d)
}

func (g *Graph) parseID(id ID) (TypeInfo, Store, error) {
	ti, ok := g.typeInfo(id.typeID())
	if !ok {
		return TypeInfo{}, nil, ErrInvalidType
	}
	if id < 1 || int(id.shardID()) > len(g.shards) {
		return TypeInfo{}, nil, fmt.Errorf("%w: invalid shard ID", ErrNotFound)
	}
	return ti, g.shards[int(id.shardID())-1], nil
}
