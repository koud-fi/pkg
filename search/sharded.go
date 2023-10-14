package search

import (
	"hash/crc64"
	"strings"

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
	shardIDs := make(map[int][]string)
	for _, id := range id {
		n := shardByID(id, len(sti.shards))
		shardIDs[n] = append(shardIDs[n], id)
	}
	var out []T
	for n, ids := range shardIDs {
		es, err := sti.shards[n].Get(ids...)
		if err != nil {
			return nil, err
		}
		out = append(out, es...)
	}
	return out, nil
}

func (sti *shardedTagIdx[T]) Query(tags []string, limit int) (QueryResult[T], error) {

	// this implementation is an extremely naive proof of concept

	var res QueryResult[T]
	for _, shard := range sti.shards {
		if err := shard.Commit(); err != nil {
			return res, err
		}
		subRes, err := shard.Query(tags, limit)
		if err != nil {
			return res, err
		}
		res.Data = append(res.Data, subRes.Data...)
		res.TotalCount += subRes.TotalCount
	}
	return res, nil
}

func (sti *shardedTagIdx[T]) Put(e ...T) {
	shardEnts := make(map[int][]T)
	for _, e := range e {
		n := shardByID(e.ID(), len(sti.shards))
		shardEnts[n] = append(shardEnts[n], e)
	}
	for n, es := range shardEnts {
		sti.shards[n].Put(es...)
	}
}

// Commit does nothing at the moment, sub-index commit is called lazily when querying
func (sti *shardedTagIdx[_]) Commit() error { return nil }

func (sti *shardedTagIdx[_]) Tags(prefix string) ([]TagInfo, error) {

	// TODO

	return []TagInfo{}, nil
}

var shardHashKeyTable = crc64.MakeTable(crc64.ISO)

func shardByID(id string, numShards int) int {
	key := id
	if i := strings.LastIndexByte(key, '.'); i > 0 {
		key = key[:i]
	}
	h := crc64.New(shardHashKeyTable)
	return int(jump.HashString(key, int32(numShards), h))
}
