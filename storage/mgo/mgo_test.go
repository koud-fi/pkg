package mgo_test

import (
	"context"
	"testing"

	"github.com/koud-fi/pkg/storage/mgo"
	"github.com/koud-fi/pkg/storage/storagetest"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Test(t *testing.T) {
	ctx := context.Background()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Disconnect(ctx)

	storagetest.TestKV(t, mgo.NewKV[string](c.Database("test").Collection("storage-test")))
}
