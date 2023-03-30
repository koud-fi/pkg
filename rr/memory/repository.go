package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/koud-fi/pkg/rr"
)

type memRW struct {
	mu     sync.RWMutex
	tables map[rr.Repository]*repoTable
}

type KeyFunc func(rr.Item) string

func NewRepository(tables map[rr.Repository]KeyFunc) rr.Writer {
	r := memRW{
		tables: make(map[rr.Repository]*repoTable, len(tables)),
	}
	for t, keyFn := range tables {
		r.tables[t] = &repoTable{
			t:     t,
			keyFn: keyFn,
			data:  make(map[string][]byte),
		}
	}
	return &r
}

/*
func (r *memRepo) Get(ctx context.Context, in rr.GetInput) (rr.GetOutput, error) {
	tables, err := resolveTables(r.tables, in)
	if err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make(rr.GetOutput, len(tables))
	for table, keys := range tables {
		items := make([]rr.Item, 0, len(keys))
		for _, k := range keys {
			if item := table.get(table.keyFn(k)); item != nil {
				items = append(items, item)
			}
		}
		out[table.t] = items
	}
	return out, nil
}
*/

func (r *memRW) Write() rr.WriteTx {
	return &writeTx{memRW: r, ops: make(map[rr.Repository][]func(*repoTable))}
}

type repoTable struct {
	t     rr.Repository
	keyFn KeyFunc
	data  map[string][]byte
}

func (rt *repoTable) get(key string) (item rr.Item) {
	if data, ok := rt.data[key]; ok {
		json.Unmarshal(data, &item) // TODO: don't ignore error
	}
	return
}

func (rt *repoTable) put(key string, item rr.Item) {
	data, _ := json.Marshal(item) // TODO: don't ignore error
	rt.data[key] = data
}

type writeTx struct {
	*memRW
	ops map[rr.Repository][]func(*repoTable)
}

func (tx *writeTx) Put(t rr.Repository, item rr.Item) {
	tx.addOp(t, func(table *repoTable) {
		table.put(table.keyFn(item), item)
	})
}

func (tx *writeTx) Update(t rr.Repository, key rr.Key, update rr.Update) {
	tx.addOp(t, func(table *repoTable) {
		var (
			k    = table.keyFn(key)
			item = table.get(k)
		)
		if item == nil {
			item = key
		}
		table.put(k, applyUpdate(item, update))
	})
}

func applyUpdate(item rr.Item, update rr.Update) rr.Item {
	for k, val := range update {
		switch v := val.(type) {
		case rr.Update:
			switch f := item[k].(type) {
			case map[string]any:
				val = applyUpdate(f, v)
			}
		}
		item[k] = val
	}
	return item
}

func (tx *writeTx) Delete(t rr.Repository, k rr.Key) {
	tx.addOp(t, func(table *repoTable) {
		delete(table.data, table.keyFn(k))
	})
}

func (tx *writeTx) addOp(t rr.Repository, op func(table *repoTable)) {
	tx.ops[t] = append(tx.ops[t], op)

}

func (tx *writeTx) Commit(ctx context.Context) error {
	tables, err := resolveTables(tx.tables, tx.ops)
	if err != nil {
		return err
	}
	tx.mu.Lock()
	defer tx.mu.Unlock()

	for table, ops := range tables {
		for _, op := range ops {
			op(table)
		}
	}
	return nil
}

func resolveTables[T any](
	tables map[rr.Repository]*repoTable, in map[rr.Repository]T,
) (map[*repoTable]T, error) {
	out := make(map[*repoTable]T, len(in))
	for t, val := range in {
		table, ok := tables[t]
		if !ok {
			return nil, fmt.Errorf("table not found: %s", t)
		}
		out[table] = val
	}
	return out, nil
}
