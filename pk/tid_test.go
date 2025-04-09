package pk_test

import (
	"testing"

	"github.com/koud-fi/pkg/pk"
)

func TestTID(t *testing.T) {
	for i := range 1000 {
		var (
			original = pk.TID(i)
			str      = original.String()
		)
		parsed, err := pk.ParseTID(str)
		if err != nil {
			t.Errorf("Failed to parse TID string %q: %v", str, err)
		}
		if parsed != original {
			t.Errorf("Parsed TID %v does not match original %v", parsed, original)
		}
	}
}
