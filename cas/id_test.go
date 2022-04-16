package cas_test

import (
	"testing"

	"github.com/koud-fi/pkg/cas"
)

func TestID(t *testing.T) {
	id := cas.NewIDFromBytes([]byte{42, 69})
	ssCid, err := cas.ParseID(id.String())
	if err != nil {
		t.Fatal(err)
	}
	hs := id.Hex()
	hsID, err := cas.ParseID(hs)
	if err != nil {
		t.Fatal(err)
	}
	if id != ssCid || id != hsID {
		t.Fatal("mismatch")
	}
	t.Log(id, hs)
}
