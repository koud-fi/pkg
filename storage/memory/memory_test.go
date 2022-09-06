package memory_test

import (
	"testing"

	"github.com/koud-fi/pkg/storage/memory"
	"github.com/koud-fi/pkg/storage/storagetest"
)

func Test(t *testing.T) {
	s := memory.NewStorage()
	storagetest.Test(t, s)

	if err := memory.Save("memstorage.temp", s, 0600); err != nil {
		t.Fatal(err)
	}
}
