package proc_test

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/koud-fi/pkg/proc"
)

type counter struct {
	N int64
}

func (c *counter) add1() int64 {
	return atomic.AddInt64(&c.N, 1)
}

func Test(t *testing.T) {
	var (
		ctx = context.Background()
		c   counter
	)
	t.Log(proc.NewOut[int64](c.add1).Invoke(ctx, nil))
	t.Log(proc.NewOut[int64](c.add1).Invoke(ctx, nil))
	t.Log(c.N)
}
