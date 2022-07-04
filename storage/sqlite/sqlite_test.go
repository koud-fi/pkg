package sqlite_test

import (
	"os"
	"testing"

	"github.com/koud-fi/pkg/storage/sqlite"
	"github.com/koud-fi/pkg/storage/storagetest"

	_ "github.com/mattn/go-sqlite3"
)

func Test(t *testing.T) {
	os.RemoveAll("temp")
	db, err := sqlite.Open("temp/test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	storagetest.Test(t, sqlite.NewStorage(db, "test"))
}
