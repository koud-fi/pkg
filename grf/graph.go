package grf

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrVersionMismatch = errors.New("version mismatch")
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

func (g *Graph) resolveAddArgs(nt NodeType) (typeID, TypeInfo, shardID, Store, error) {
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
	return ti, g.shardForID(id.shardID()), nil
}

func (g *Graph) shardForID(id shardID) Store {
	return g.shards[int(id)-1]
}

func (g *Graph) parseEdgeArgs(from ID, et EdgeType) (TypeInfo, EdgeTypeID, Store, error) {
	ti, s, err := g.parseID(from)
	if err != nil {
		return TypeInfo{}, 0, nil, err
	}
	etID, ok := ti.edgeTypeMap[et]
	if !ok {
		return TypeInfo{}, 0, nil, ErrInvalidEdgeType
	}
	return ti, etID, s, err
}
