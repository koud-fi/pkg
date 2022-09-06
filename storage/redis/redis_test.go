package redis_test

import (
	"testing"

	"github.com/koud-fi/pkg/storage/redis"
	"github.com/koud-fi/pkg/storage/storagetest"
)

func Test(t *testing.T) {
	c := redis.Open(&redis.Options{Addr: "localhost:6379"})
	defer c.Close()

	storagetest.Test(t, redis.NewStorage(c, redis.KeyPrefix("__TEST:")))
}
