package grf

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidType   = errors.New("invalid type")
	ErrAlreadyExists = errors.New("already exists")
)

type Graph struct {
	m       Mapper
	shards  []Store
	counter int64
	types   map[typeID]NodeType
	typeIDs map[NodeType]typeID
}

func New(m Mapper, s ...Store) *Graph {
	if len(s) == 0 {
		panic("no stores")
	}
	return &Graph{
		shards:  s,
		types:   make(map[typeID]NodeType),
		typeIDs: make(map[NodeType]typeID),
	}
}

func (g *Graph) Register(nt NodeType, id typeID) {
	if nt == "" || id < 1 {
		panic("invalid type/ID") // TODO: better type validation
	}

	// TODO: prevent duplicate types/ids

	g.types[id] = nt
	g.typeIDs[nt] = id
}

func (g *Graph) Node(id ID) (*Node, error) {
	nt, s, err := g.parseID(id)
	if err != nil {
		return nil, err
	}
	ns, err := s.Node(nt, id.localID())
	if err != nil {
		return nil, err
	}
	if len(ns) == 0 {
		return nil, fmt.Errorf("%w: %d", ErrNotFound, id)
	}
	return &Node{
		id: id,
		t:  nt,
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
	typeID, ok := g.typeIDs[nt]
	if !ok {
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
	nt, s, err := g.parseID(id)
	if err != nil {
		return err
	}
	return s.UpdateNode(nt, id.localID(), marshal(v))
}

func (g *Graph) DeleteNode(id ID) error {
	nt, s, err := g.parseID(id)
	if err != nil {
		return err
	}

	// TODO: support batch delete

	return s.DeleteNode(nt, id.localID())
}

func (g *Graph) parseID(id ID) (NodeType, Store, error) {
	nt, ok := g.types[id.typeID()]
	if !ok {
		return "", nil, ErrInvalidType
	}
	if id < 1 || int(id.shardID()) > len(g.shards) {
		return "", nil, fmt.Errorf("%w: invalid shard ID", ErrNotFound)
	}
	return nt, g.shards[int(id.shardID())-1], nil
}

func marshal(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
