package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type KV[T any] struct {
	coll *mongo.Collection
}

func NewKV[T any](c *mongo.Collection) *KV[T] {
	return &KV[T]{coll: c}
}

func (kv *KV[T]) Get(ctx context.Context, key string) (T, error) {
	var d doc[T]
	err := kv.coll.FindOne(ctx, bson.M{"_id": key}).Decode(&d)
	return d.Data, err
}

func (kv *KV[T]) Put(ctx context.Context, key string, value T) error {
	_, err := kv.coll.UpdateOne(ctx, bson.M{
		"_id": key,
	}, bson.M{
		"$set": doc[T]{
			Data:      value,
			UpdatedAt: time.Now(),
		},
	}, options.Update().SetUpsert(true))
	return err
}

func (kv *KV[T]) Delete(ctx context.Context, keys ...string) error {
	_, err := kv.coll.DeleteMany(ctx, bson.M{"_id": bson.M{
		"$in": keys,
	}})
	return err
}

type doc[T any] struct {
	ID        string    `bson:"_id"`
	Data      T         `bson:"data"`
	UpdatedAt time.Time `bson:"updated_at"`
}
