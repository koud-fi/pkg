package mgo

import (
	"context"
	"fmt"

	"github.com/koud-fi/pkg/proto/rr"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRW struct {
	repos map[rr.Repository]*mgoRepo
}

func NewRW(db *mongo.Database, keys map[rr.Repository][]string) rr.ReadWriter {
	repos := make(map[rr.Repository]*mgoRepo)
	for r := range keys {
		repos[r] = &mgoRepo{
			r:    r,
			coll: db.Collection(string(r)),

			// TODO: key/id funcs

		}
	}
	return &mongoRW{repos: repos}
}

func (m *mongoRW) Read() rr.ReadTx {
	return &readTx{mongoRW: m, keys: make(map[rr.Repository][]rr.Item)}
}

type readTx struct {
	*mongoRW
	keys map[rr.Repository][]rr.Key
}

func (tx *readTx) Get(r rr.Repository, keys ...rr.Key) {
	tx.keys[r] = append(tx.keys[r], keys...)
}

func (tx *readTx) Execute(ctx context.Context) (map[rr.Repository][]rr.Item, error) {
	/*
		res := make(map[rr.Repository][]rr.Item, len(tx.keys))
		for r, keys := range tx.keys {
			cur, err := tx.db.Collection(string(r)).Find(ctx, bson.M{
				"_id": bson.M{"$in": keys},
			})
			if err != nil {
				return nil, err
			}
			var items []rr.Item
			if err := cur.All(ctx, &items); err != nil {
				return nil, err
			}
			res[r] = items
		}
		return res, nil
	*/
	panic("TODO")
}

func (m *mongoRW) Write() rr.WriteTx {
	return &writeTx{mongoRW: m, models: make(map[rr.Repository][]mongo.WriteModel)}
}

type writeTx struct {
	*mongoRW
	models map[rr.Repository][]mongo.WriteModel
}

func (tx *writeTx) Put(r rr.Repository, item rr.Item) {

	// TODO: support upsertion (once such option exists in the writer interface)

	tx.models[r] = append(tx.models[r], mongo.NewInsertOneModel().SetDocument(item))
}

func (tx *writeTx) Update(r rr.Repository, key rr.Key, update rr.Update) {
	set := make(bson.M, len(update))
	for k, v := range update {
		set[k] = v
	}
	if len(set) == 0 {
		return
	}
	up := bson.M{"$set": set}
	tx.models[r] = append(tx.models[r], mongo.NewUpdateOneModel().SetFilter(key).SetUpdate(up))
}

func (tx *writeTx) Delete(r rr.Repository, key rr.Key) {
	tx.models[r] = append(tx.models[r], mongo.NewDeleteOneModel().SetFilter(key))
}

func (tx *writeTx) Commit(ctx context.Context) error {

	// TODO: look into possible multi collection transactions?

	/*
		for r, models := range tx.models {
			if _, err := tx.db.Collection(string(r)).BulkWrite(ctx, models); err != nil {
				return err
			}
		}
		return nil
	*/
	panic("TODO")
}

type mgoRepo struct {
	r     rr.Repository
	idFn  func(rr.Key) bson.M
	keyFn func(bson.M) rr.Key
	coll  *mongo.Collection
}

func resolveRepos[T any](
	repos map[rr.Repository]*mgoRepo, in map[rr.Repository]T,
) (map[*mgoRepo]T, error) {
	out := make(map[*mgoRepo]T, len(in))
	for r, val := range in {
		repo, ok := repos[r]
		if !ok {
			return nil, fmt.Errorf("repository not found: %s", r)
		}
		out[repo] = val
	}
	return out, nil
}
