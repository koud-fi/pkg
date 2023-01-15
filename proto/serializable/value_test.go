package serializable_test

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/koud-fi/pkg/proto/serializable"
)

type Thing struct {
	A string
	B string
}

func (t Thing) Encode(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%s:%s", t.A, t.B)
	return err
}

func (t Thing) Decode(w io.Reader) error {

	// ???

	panic("TODO")
}

func TestValue(t *testing.T) {
	var data struct {
		T1, T2 serializable.Value[Thing]
	}
	data.T1.Data.A = "A1"
	data.T1.Data.B = "B1"
	data.T2.Data.A = "A2"
	data.T2.Data.B = "B2"

	buf, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(buf)) // TODO: proper validation
}
