package diriter_test

import (
	"os"
	"testing"

	"github.com/koud-fi/pkg/rx"
	"github.com/koud-fi/pkg/rx/diriter"
)

func TestDirIter(t *testing.T) {
	t.Log(rx.Drain(rx.Log(diriter.Paths(diriter.New(os.DirFS(".."), ".")), "")))
}
