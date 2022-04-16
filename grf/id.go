package grf

const (
	localIDMax    = 1<<31 - 1
	typeIDMax     = 1<<15 - 1
	typeIDOffset  = 31
	shardIDMax    = 1<<15 - 1
	shardIDOffset = 46
)

type ID int64
type shardID int16
type typeID int16

func newID(s shardID, t typeID, l LocalID) ID {
	return ID(int64(t)<<typeIDOffset | int64(s)<<shardIDOffset | int64(l))
}

func (id ID) typeID() typeID   { return typeID(id >> typeIDOffset & typeIDMax) }
func (id ID) shardID() shardID { return shardID(id >> shardIDOffset & shardIDMax) }
func (id ID) localID() LocalID { return LocalID(id & localIDMax) }
