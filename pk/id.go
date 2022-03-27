package pk

import "fmt"

const (
	typeIDMax     = 0x3FF
	typeIDOffset  = 36
	shardIDMax    = 0xFFFF
	shardIDOffset = 46
	localIDMax    = 0xFFFFFFFFF
)

// ID is a 64 bit ID that is a combination of 10 bit non-zero type ID, 16 bit shard ID
// and a 36 bit non-zero "local ID" inside the shard. One bit is saved for the future
// and the last bit is not used, so it doesn't matter if you use
// signed or unsigned 64 bit integer type to store the value.
type ID int64

// NewID returns a new ID using the given type, shard and local ID,
// it will panic if any of the given IDs are out of range.
func NewID(typeID int16, shardID uint16, localID int64) ID {
	if typeID < 1 || typeID > typeIDMax {
		panic(fmt.Sprintf("pk: type ID out of range (1-%d): %d", typeIDMax, typeID))
	}
	if localID < 1 || localID > typeIDMax {
		panic(fmt.Sprintf("pk: local ID out of range (1-%d): %d", localIDMax, localID))
	}
	return ID(int64(typeID)<<typeIDOffset | int64(shardID)<<shardIDOffset | localID)
}

// TypeID returns the underlying 10 bit type ID.
func (t ID) TypeID() int16 { return int16(t >> typeIDOffset & typeIDMax) }

// ShardID returns the underlying 16 bit shard ID.
func (t ID) ShardID() uint16 { return uint16(t >> shardIDOffset & shardIDMax) }

// LocalID returns the underlying 36 bit local ID inside the shard.
func (t ID) LocalID() int64 { return int64(t & localIDMax) }
