package sqlgrf_test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/koud-fi/pkg/proto/grf/grftest"
	"github.com/koud-fi/pkg/proto/grf/sqlgrf"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "temp/test.db"

func Test(t *testing.T) {
	dir := filepath.Dir(dbPath)
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared", dbPath))
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	grftest.Test(t, sqlgrf.NewStore(db))
}
