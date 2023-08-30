package memory_test

import (
	"fmt"
	"testing"

	"github.com/koud-fi/pkg/proto/rr"
	"github.com/koud-fi/pkg/proto/rr/memory"
	"github.com/koud-fi/pkg/proto/rr/rrtest"
)

func TestRW(t *testing.T) {
	rrtest.Run(t, memory.NewRW(map[rr.Repository]memory.KeyFunc{
		rrtest.Repository: func(item rr.Item) string { return fmt.Sprint(item["id"]) },
	}))
}
