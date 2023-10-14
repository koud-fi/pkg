package search

import (
	"hash/crc64"

	"github.com/koud-fi/pkg/jump"
)

type shardedTagIdx[T Entry] struct {
	shards []TagIndex[T]
}

func NewShardedTagIndex[T Entry](numShards int, shardInitFn func(n int) TagIndex[T]) TagIndex[T] {
	shards := make([]TagIndex[T], numShards)
	for n := 0; n < numShards; n++ {
		shards[n] = shardInitFn(n)
	}
	return &shardedTagIdx[T]{shards: shards}
}

func (sti *shardedTagIdx[T]) Get(id ...string) ([]T, error) {

	// ???

	panic("TODO")
}

func (sti *shardedTagIdx[T]) Query(tags []string, limit int) (QueryResult[T], error) {

	// ???

	panic("TODO")
}

func (sti *shardedTagIdx[T]) Put(e ...T) {

	// ???

	panic("TODO")
}

func (sti *shardedTagIdx[_]) Commit() error {

	// ???

	panic("TODO")
}

func (sti *shardedTagIdx[_]) Tags(prefix string) ([]TagInfo, error) {

	// TODO

	return []TagInfo{}, nil
}

var shardHashKeyTable = crc64.MakeTable(crc64.ISO)

func shardByID(id string, numShards int) int {
	h := crc64.New(shardHashKeyTable)
	return int(jump.HashString(id, int32(numShards), h))
}
