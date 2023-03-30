package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/koud-fi/pkg/rr"
)

type memRW struct {
	mu    sync.RWMutex
	repos map[rr.Repository]*memRepo
}

type KeyFunc func(rr.Item) string

func NewRW(repos map[rr.Repository]KeyFunc) rr.Writer {
	m := memRW{
		repos: make(map[rr.Repository]*memRepo, len(repos)),
	}
	for r, keyFn := range repos {
		m.repos[r] = &memRepo{
			r:     r,
			keyFn: keyFn,
			data:  make(map[string][]byte),
		}
	}
	return &m
}

func (m *memRW) Read() rr.ReadTx {
	return &readTx{memRW: m}
}

type readTx struct {
	req map[rr.Repository][]rr.Key
	*memRW
}

func (tx *readTx) Get(r rr.Repository, key rr.Key) {
	tx.req[r] = append(tx.req[r], key)
}

func (tx *readTx) Execute(ctx context.Context) (map[rr.Repository][]rr.Item, error) {
	repos, err := resolveRepos(tx.repos, tx.req)
	if err != nil {
		return nil, err
	}
	tx.mu.RLock()
	defer tx.mu.RUnlock()

	out := make(map[rr.Repository][]rr.Item, len(repos))
	for repo, keys := range repos {
		items := make([]rr.Item, 0, len(keys))
		for _, k := range keys {
			if item := repo.get(repo.keyFn(k)); item != nil {
				items = append(items, item)
			}
		}
		out[repo.r] = items
	}
	return out, nil
}

func (m *memRW) Write() rr.WriteTx {
	return &writeTx{memRW: m, ops: make(map[rr.Repository][]func(*memRepo))}
}

type writeTx struct {
	*memRW
	ops map[rr.Repository][]func(*memRepo)
}

func (tx *writeTx) Put(r rr.Repository, item rr.Item) {
	tx.addOp(r, func(repo *memRepo) {
		repo.put(repo.keyFn(item), item)
	})
}

func (tx *writeTx) Update(r rr.Repository, key rr.Key, update rr.Update) {
	tx.addOp(r, func(repo *memRepo) {
		var (
			k    = repo.keyFn(key)
			item = repo.get(k)
		)
		if item == nil {
			item = key
		}
		repo.put(k, applyUpdate(item, update))
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

func (tx *writeTx) Delete(r rr.Repository, k rr.Key) {
	tx.addOp(r, func(repo *memRepo) {
		delete(repo.data, repo.keyFn(k))
	})
}

func (tx *writeTx) addOp(r rr.Repository, op func(repo *memRepo)) {
	tx.ops[r] = append(tx.ops[r], op)

}

func (tx *writeTx) Commit(ctx context.Context) error {
	repos, err := resolveRepos(tx.repos, tx.ops)
	if err != nil {
		return err
	}
	tx.mu.Lock()
	defer tx.mu.Unlock()

	for repo, ops := range repos {
		for _, op := range ops {
			op(repo)
		}
	}
	return nil
}

type memRepo struct {
	r     rr.Repository
	keyFn KeyFunc
	data  map[string][]byte
}

func (mr *memRepo) get(key string) (item rr.Item) {
	if data, ok := mr.data[key]; ok {
		json.Unmarshal(data, &item) // TODO: don't ignore error
	}
	return
}

func (mr *memRepo) put(key string, item rr.Item) {
	data, _ := json.Marshal(item) // TODO: don't ignore error
	mr.data[key] = data
}

func resolveRepos[T any](
	repos map[rr.Repository]*memRepo, in map[rr.Repository]T,
) (map[*memRepo]T, error) {
	out := make(map[*memRepo]T, len(in))
	for r, val := range in {
		repo, ok := repos[r]
		if !ok {
			return nil, fmt.Errorf("repository not found: %s", r)
		}
		out[repo] = val
	}
	return out, nil
}
