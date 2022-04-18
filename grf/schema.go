package grf

import "fmt"

type schema struct {
	types   []TypeInfo
	typeMap map[NodeType]typeID
}

type TypeInfo struct {
	Type  NodeType
	Edges []EdgeTypeInfo
}

type EdgeTypeInfo struct {
	Type EdgeType
}

func (s *schema) register(ti TypeInfo) {
	if _, ok := s.typeMap[ti.Type]; ok {
		panic(fmt.Sprintf("duplicate type: %s", ti.Type))
	}
	s.types = append(s.types, ti)
	s.typeMap[ti.Type] = typeID(len(s.types))

	// TODO: handle edges

}

func (s schema) typeInfo(id typeID) (TypeInfo, bool) {
	n := int(id) - 1
	if n < 0 || n >= len(s.types) {
		return TypeInfo{}, false
	}
	return s.types[n], true
}
