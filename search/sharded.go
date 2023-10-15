package search

import (
	"hash/crc64"
	"math/rand"
	"strings"
	"sync"

	"github.com/koud-fi/pkg/jump"
)

type ShardedTagIndex[T Entry] struct {
	shards     []TagIndex[T]
	queryOrder []int
	resPool    *sync.Pool
}

func NewShardedTagIndex[T Entry](
	numShards int32, shardInitFn func(n int32) TagIndex[T],
) ShardedTagIndex[T] {
	shards := make([]TagIndex[T], numShards)
	for n := int32(0); n < numShards; n++ {
		shards[n] = shardInitFn(n)
	}
	return ShardedTagIndex[T]{
		shards:     shards,
		queryOrder: queryOrder(len(shards), 0),
		resPool: &sync.Pool{
			New: func() any { return &QueryResult[T]{Data: make([]T, 0, 1<<10)} },
		},
	}
}

func (sti ShardedTagIndex[T]) WithSeed(seed int64) ShardedTagIndex[T] {
	sti.queryOrder = queryOrder(len(sti.shards), seed)
	return sti
}

func queryOrder(numShards int, seed int64) []int {
	orders := make([]int, numShards)
	for i := 0; i < numShards; i++ {
		orders[i] = i
	}
	if seed != 0 {
		rand.New(rand.NewSource(seed)).Shuffle(len(orders), func(i, j int) {
			orders[i], orders[j] = orders[j], orders[i]
		})
	}
	return orders
}

func (sti ShardedTagIndex[T]) Get(id ...string) ([]T, error) {
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

func (sti ShardedTagIndex[T]) Query(dst *QueryResult[T], tags []string, limit int) error {
	dst.Reset()

	// TODO: stop if total limit is reached
	// TODO: do another pass if total limit is not reached and shards hit their individual limits
	// TODO: do some form of magic for better ordering, don't blindly sort to retain seed variance

	preAlloc := limit/8 + 1
	if cap(dst.Data) < preAlloc/2 {
		dst.Data = append(dst.Data, make([]T, 0, preAlloc-cap(dst.Data))...)
	}
	subLimit := limit / len(sti.shards) * 4

	for _, i := range sti.queryOrder {
		shard := sti.shards[i]
		if err := shard.Commit(); err != nil {
			return err
		}
		var (
			res = sti.resPool.Get().(*QueryResult[T])
			err = shard.Query(res, tags, subLimit)
		)
		dst.Data = append(dst.Data, res.Data...)
		dst.TotalCount += res.TotalCount

		sti.resPool.Put(res)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sti ShardedTagIndex[T]) Put(e ...T) {
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
func (sti ShardedTagIndex[_]) Commit() error { return nil }

func (sti ShardedTagIndex[_]) Tags(prefix string) ([]TagInfo, error) {

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
