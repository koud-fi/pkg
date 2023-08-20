package bolt_test

import (
	"os"
	"testing"

	"github.com/koud-fi/pkg/storage/bolt"
	"github.com/koud-fi/pkg/storage/storagetest"
)

func Test(t *testing.T) {
	os.RemoveAll("temp")
	db, err := bolt.Open("temp/test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	storagetest.Test(t, bolt.NewStorage(db, "test-bucket"))
}
