package mgo

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/koud-fi/pkg/rx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type doc[T any] struct {
	ID        string    `bson:"_id,omitempty"`
	Data      T         `bson:"data"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type KV[T any] struct {
	coll *mongo.Collection
}

func NewKV[T any](c *mongo.Collection) *KV[T] {
	return &KV[T]{coll: c}
}

func (kv *KV[T]) Get(ctx context.Context, key string) (T, error) {
	var d doc[T]
	err := kv.coll.FindOne(ctx, bson.M{"_id": key}).Decode(&d)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		err = os.ErrNotExist
	}
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

func (kv *KV[T]) Iter(ctx context.Context, after string) rx.Iter[rx.Pair[string, T]] {
	return &iter[T]{ctx: ctx, coll: kv.coll, after: after}
}

type iter[T any] struct {
	ctx   context.Context
	coll  *mongo.Collection
	after string

	cur *mongo.Cursor
	err error
	d   doc[T]
}

func (it *iter[T]) Next() bool {
	if it.cur == nil {
		filter := make(bson.M)
		if it.after != "" {
			filter["_id"] = bson.M{"$lt": it.after}
		}
		opts := options.Find().SetSort(bson.M{"_id": 1})
		it.cur, it.err = it.coll.Find(it.ctx, filter, opts)
	}
	if it.err != nil {
		return false
	}
	if !it.cur.Next(it.ctx) {
		return false
	}
	it.err = it.cur.Decode(&it.d)
	return it.err == nil
}

func (it *iter[T]) Value() rx.Pair[string, T] {
	return rx.NewPair[string, T](it.d.ID, it.d.Data)
}

func (it *iter[_]) Close() error {
	if err := it.cur.Close(it.ctx); err != nil {
		return err
	}
	return it.err
}
