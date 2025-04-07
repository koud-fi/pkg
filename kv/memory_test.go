package kv_test

import (
	"testing"

	"github.com/koud-fi/pkg/kv"
	"github.com/koud-fi/pkg/kv/kvtest"
)

func TestMemoryStorage(t *testing.T) {
	if err := kvtest.TestKV(kv.NewMemoryStorage[string, string]()); err != nil {
		t.Fatal(err)
	}
}
