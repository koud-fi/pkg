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

func (g *Graph) Node(id ID) *Node {
	n := Node{id: id}
	if n.ti, n.s, n.err = g.parseID(id); n.err != nil {
		return &n
	}
	var ns []NodeData
	if ns, n.err = n.s.Node(n.ti.Type, id.localID()); n.err != nil {
		return &n
	}
	if len(ns) == 0 {
		n.err = fmt.Errorf("%w: %d", ErrNotFound, id)
		return &n
	}
	n.d = ns[0]
	return &n
}

func (g *Graph) MappedNode(nt NodeType, key string, add bool) *Node {
	id, err := g.m.Map(nt, key)
	if err != nil {
		if add && err == ErrNotFound {
			n, err := g.AddNode(nt, nil)
			if err != nil {
				return &Node{err: err}
			}
			n.err = g.m.SetMapping(nt, key, n.id)
			return n
		}
		return &Node{err: err}
	}
	return g.Node(id)
}

func (g *Graph) AddNode(nt NodeType, v any) (*Node, error) {
	typeID, ti, shardID, s, err := g.resolveAddParams(nt)
	if err != nil {
		return nil, err
	}
	data := marshal(v)
	localID, ts, err := s.AddNode(nt, data)
	if err != nil {
		return nil, err
	}
	return &Node{
		s:  s,
		id: newID(shardID, typeID, localID),
		d: NodeData{
			ID:        localID,
			Data:      data,
			Timestamp: ts,
		},
		ti: ti,
	}, nil
}

func (g *Graph) DeleteNode(id ID) error {
	ti, s, err := g.parseID(id)
	if err != nil {
		return err
	}
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
	return s.SetEdge(ti.Type, e.d)
}

func (g *Graph) resolveAddParams(nt NodeType) (typeID, TypeInfo, shardID, Store, error) {
	typeID := g.typeMap[nt]
	ti, ok := g.typeInfo(typeID)
	if !ok {
		return 0, TypeInfo{}, 0, nil, ErrInvalidType
	}
	var (
		shardNum = int(atomic.AddInt64(&g.counter, 1) % int64(len(g.shards)))
		s        = g.shards[shardNum]
		shardID  = shardID(shardNum + 1)
	)
	return typeID, ti, shardID, s, nil
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
