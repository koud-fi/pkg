package pk_test

import (
	"encoding/json"
	"testing"

	"github.com/koud-fi/pkg/pk"
)

func TestTID(t *testing.T) {
	for i := range 1000 {
		var original pk.TID
		original.Set(int64(i))

		str := original.String()

		parsed, err := pk.ParseTID(str)
		if err != nil {
			t.Errorf("Failed to parse TID string %q: %v", str, err)
		}
		if parsed != original {
			t.Errorf("Parsed TID %v does not match original %v", parsed, original)
		}
	}
}

func TestTID_MarshalJSON(t *testing.T) {
	type TIDAlias pk.TID
	type TIDPtrAlias *pk.TID
	type TIDWrapper struct{ pk.TID }

	var tid pk.TID
	tid.Set(123)

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
