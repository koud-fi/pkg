package pk_test

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/koud-fi/pkg/pk"
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
