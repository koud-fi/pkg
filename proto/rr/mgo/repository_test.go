package mgo_test

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/koud-fi/pkg/rr"
	"github.com/koud-fi/pkg/rr/mgo"
	"github.com/koud-fi/pkg/rr/rrtest"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var connURL = flag.String("url", "mongodb://localhost:27017", "MongoDB connection URL")

func TestRW(t *testing.T) {
	client, err := mongo.NewClient(options.Client().ApplyURI(*connURL))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Fatal(err)
	}
	rrtest.Run(t, mgo.NewRW(client.Database("test"), map[rr.Repository][]string{
		rrtest.Repository: {"id"},
	}))
}
