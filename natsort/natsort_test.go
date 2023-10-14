package natsort_test

import (
	"reflect"
	"testing"

	"github.com/koud-fi/pkg/natsort"
)

func TestSort(t *testing.T) {
	data := []string{
		"a1.txt",
		"a10.txt",
		"a200.txt",
		"a100.txt",
		"a2.txt",
		"a20.txt",
		"a3 xyz.txt",
		"a4.txt",
		"a4_v3.txt",
		"a4_v2.txt",
		"a40.txt",
		"a9.txt",
		"a0.txt",
		"aa.txt",
	}
	dataOK := []string{
		"a0.txt",
		"a1.txt",
		"a2.txt",
		"a3 xyz.txt",
		"a4.txt",
		"a4_v2.txt",
		"a4_v3.txt",
		"a9.txt",
		"a10.txt",
		"a20.txt",
		"a40.txt",
		"a100.txt",
		"a200.txt",
		"aa.txt",
	}
	natsort.Strings(data)
	if !reflect.DeepEqual(data, dataOK) {
		t.Fatalf("data mismatch\n\texpected: %v\n\tgot:      %v", dataOK, data)
	}
}
