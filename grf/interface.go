package grf

import "time"

type NodeData struct {
	ID        LocalID   `json:"id"`
	Data      []byte    `json:"data"`
	Timestamp time.Time `json:"ts"`
}

type LocalID int32

type EdgeData struct {
	From     LocalID  `json:"from"`
	Type     EdgeType `json:"type"`
	To       ID       `json:"to"`
	Sequence int64    `json:"seq,omitempty"`
	Data     []byte   `json:"data,omitempty"`
}

type Store interface {
	Node(nt NodeType, id ...LocalID) ([]NodeData, error)
	NodeRange(nt NodeType, after LocalID, limit int) ([]NodeData, error)

	Edge(nt NodeType, from LocalID, et EdgeType, to ...ID) ([]EdgeData, error)
	EdgeCount(nt NodeType, from LocalID, et ...EdgeType) (map[EdgeType]int, error)
	EdgeRange(nt NodeType, from LocalID, et EdgeType, offset, limit int) ([]EdgeData, error)
	// TODO: sequence based edge range method

	AddNode(nt NodeType, data []byte) (LocalID, time.Time, error)
	UpdateNode(nt NodeType, id LocalID, data []byte) error
	DeleteNode(nt NodeType, id ...LocalID) error

	SetEdge(nt NodeType, e ...EdgeData) error
	// TODO: edge type changing method
	DeleteEdge(nt NodeType, from LocalID, et EdgeType, to ...ID) error
}

type Mapper interface {
	Map(nt NodeType, key string) (ID, error)

	SetMapping(nt NodeType, key string, id ID) error
	DeleteMapping(nt NodeType, key ...string) error
}
