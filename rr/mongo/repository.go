package mongo

import (
	"context"

	"github.com/koud-fi/pkg/rr"

	"go.mongodb.org/mongo-driver/mongo"
)

type writeTx struct {
	db     *mongo.Database
	models map[rr.Repository][]mongo.WriteModel
}

func (tx *writeTx) Put(r rr.Repository, item rr.Item) {
	tx.models[r] = append(tx.models[r], mongo.NewInsertOneModel().SetDocument(item))
}

func (tx *writeTx) Update(r rr.Repository, key rr.Key, update rr.Update) {
	ops := update

	// TODO: create correct update ops

	tx.models[r] = append(tx.models[r], mongo.NewUpdateOneModel().SetFilter(key).SetUpdate(ops))
}

func (tx *writeTx) Delete(r rr.Repository, key rr.Key) {
	tx.models[r] = append(tx.models[r], mongo.NewDeleteOneModel().SetFilter(key))
}

func (tx *writeTx) Commit(ctx context.Context) error {

	// TODO: look into possible multi collection transactions?

	for r, models := range tx.models {
		if _, err := tx.db.Collection(string(r)).BulkWrite(ctx, models); err != nil {
			return err
		}
	}
	return nil
}
