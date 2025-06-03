package pk_test

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/koud-fi/pkg/pk"
	_ "github.com/mattn/go-sqlite3"
)

func TestTID(t *testing.T) {
	for i := range 1 << 20 {
		testTID(t, int64(i))
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for range 1 << 20 {
		testTID(t, rnd.Int63())
	}
}

func testTID(t *testing.T, n int64) {
	var (
		original = pk.NewRawValueTID(n)
		raw      = original.RawValue()
		str      = original.String()
	)
	parsed, err := pk.ParseTID(str)
	if err != nil {
		t.Errorf("Failed to parse TID string %q: %v", str, err)
	}
	if parsed != original {
		t.Errorf("Parsed TID %v does not match original %v", parsed, original)
	}
	if parsed.RawValue() != raw {
		t.Errorf("Parsed TID value %v does not match raw value %v", parsed.RawValue(), raw)
	}
	if raw != n {
		t.Errorf("Raw value %d does not match expected value %d", raw, n)
	}
}

func TestTID_MarshalJSON(t *testing.T) {
	type TIDAlias pk.TID
	type TIDPtrAlias *pk.TID
	type TIDWrapper struct{ pk.TID }

	tid := pk.NewRawValueTID(123)

	testTIDMarshalJSON(t, tid, `"ew"`)
	testTIDMarshalJSON(t, TIDAlias(tid), "{}") // NOTE: Non-pointer aliasing breaks MarshalJSON
	testTIDMarshalJSON(t, TIDPtrAlias(&tid), `"ew"`)
	testTIDMarshalJSON(t, TIDWrapper{tid}, `"ew"`)
}

func testTIDMarshalJSON(t *testing.T, tid any, expected string) {
	json, err := json.Marshal(tid)
	if err != nil {
		t.Errorf("Failed to marshal TID: %v", err)
	}
	if string(json) != expected {
		t.Errorf("Unexpected JSON output: %s, expected %s", json, expected)
	}
}

func TestTID_DatabaseRoundTrip(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create test table
	_, err = db.Exec("CREATE TABLE test_tids (id INTEGER PRIMARY KEY, tid_value INTEGER)")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	// Test cases: focus on basic functionality
	testCases := []struct {
		name string
		tid  pk.TID
	}{
		{"regular_tid", pk.NewTID(time.Now(), 1, 42, false)},
		{"serial_tid", pk.NewSerialTID(1, 456)},
		{"raw_value_tid", pk.NewRawValueTID(999999)},
		{"zero_tid", pk.TID{}},
	}
	// Test round-trip for each TID
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Insert TID into database
			_, err := db.Exec("INSERT INTO test_tids (tid_value) VALUES (?)", tc.tid)
			if err != nil {
				t.Fatalf("Failed to insert TID: %v", err)
			}
			// Read TID back from database
			var scannedTID pk.TID
			err = db.QueryRow("SELECT tid_value FROM test_tids ORDER BY id DESC LIMIT 1").Scan(&scannedTID)
			if err != nil {
				t.Fatalf("Failed to scan TID: %v", err)
			}
			// Compare the scanned TID with the original
			if scannedTID.RawValue() != tc.tid.RawValue() {
				t.Errorf("Scanned TID raw value %d does not match original %d",
					scannedTID.RawValue(), tc.tid.RawValue())
			}
			// Verify that the TID can be converted back to string correctly
			if !tc.tid.IsZero() && scannedTID.String() != tc.tid.String() {
				t.Errorf("String representation mismatch: got %q, want %q",
					scannedTID.String(), tc.tid.String())
			}
			// Clean up for next test
			_, err = db.Exec("DELETE FROM test_tids")
			if err != nil {
				t.Fatalf("Failed to clean up: %v", err)
			}
		})
	}
}
