package memory_test

import (
	"testing"

	"github.com/koud-fi/pkg/storage/memory"
	"github.com/koud-fi/pkg/storage/storagetest"
)

func Test(t *testing.T) {
	s := memory.NewStorage()
	storagetest.Test(t, s)
}
