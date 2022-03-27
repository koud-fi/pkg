package localdisk_test

import (
	"os"
	"testing"

	"github.com/koud-fi/pkg/storage/localdisk"
	"github.com/koud-fi/pkg/storage/storagetest"
)

func Test(t *testing.T) {
	os.RemoveAll("temp")
	s, err := localdisk.NewStorage("temp")
	if err != nil {
		t.Fatal(err)
	}
	storagetest.Test(t, s)
}
