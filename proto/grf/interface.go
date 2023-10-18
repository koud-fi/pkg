package grf

type NodeData struct {
	ID      LocalID
	Data    []byte
	Version int64
}

type LocalID int32

type EdgeData struct {
	From     LocalID
	TypeID   EdgeTypeID
	To       ID
	Sequence int64
	Data     []byte
}

type EdgeTypeID int32

type EdgeInfoData struct {
	TypeID EdgeTypeID
	EdgeInfo
}

type Store interface {
	Node(nt NodeType, id ...LocalID) ([]NodeData, error)
	NodeRange(nt NodeType, after LocalID, limit int) ([]NodeData, error)

	Edge(nt NodeType, from LocalID, et EdgeTypeID, to ...ID) ([]EdgeData, error)
	EdgeInfo(nt NodeType, from LocalID, et ...EdgeTypeID) ([]EdgeInfoData, error)
	EdgeRange(nt NodeType, from LocalID, et EdgeTypeID, offset, limit int) ([]EdgeData, error)
	// TODO: sequence based edge range method

	AddNode(nt NodeType, data []byte) (LocalID, int64, error)
	UpdateNode(nt NodeType, id LocalID, data []byte, currentVersion int64) error
	DeleteNode(nt NodeType, id ...LocalID) error

	SetEdge(nt NodeType, e ...EdgeData) error
	// TODO: edge type changing method
	DeleteEdge(nt NodeType, from LocalID, et EdgeTypeID, to ...ID) error
}

type Mapper interface {
	Map(nt NodeType, key string) (ID, error)

	SetMapping(nt NodeType, key string, id ID) error
	DeleteMapping(nt NodeType, key ...string) error
}
