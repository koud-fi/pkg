package grf

import "fmt"

type schema struct {
	types   []TypeInfo
	typeMap map[NodeType]typeID
}

type TypeInfo struct {
	Type  NodeType
	Edges []EdgeTypeInfo

	edgeTypeMap map[EdgeType]EdgeTypeID
}

type EdgeTypeInfo struct {
	Type EdgeType
}

func (s *schema) register(ti TypeInfo) {
	if _, ok := s.typeMap[ti.Type]; ok {
		panic(fmt.Sprintf("duplicate type: %s", ti.Type))
	}
	ti.edgeTypeMap = make(map[EdgeType]EdgeTypeID, len(ti.Edges))
	for i, e := range ti.Edges {
		if _, ok := ti.edgeTypeMap[e.Type]; ok {
			panic(fmt.Sprintf("duplicate edge type: %s on type: %s", e.Type, ti.Type))
		}
		ti.edgeTypeMap[e.Type] = EdgeTypeID(i + 1)
	}
	s.types = append(s.types, ti)
	s.typeMap[ti.Type] = typeID(len(s.types))
}

func (s schema) typeInfo(id typeID) (TypeInfo, bool) {
	n := int(id) - 1
	if n < 0 || n >= len(s.types) {
		return TypeInfo{}, false
	}
	return s.types[n], true
}
