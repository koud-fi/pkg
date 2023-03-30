package rrtest

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/koud-fi/pkg/rr"
)

const Repository = "test"

func Run(t *testing.T, rw rr.ReadWriter) {
	ctx := context.Background()

	wtx1 := rw.Write()
	wtx1.Put(Repository, rr.Item{
		"id":    1,
		"value": "Hello?",
	})
	wtx1.Put(Repository, rr.Item{
		"id":    2,
		"value": "World!!!", "n": 42,
		"list": []int64{1, 2, 3},
	})
	wtx1.Put(Repository, rr.Item{
		"id": 10,
		"stuff": map[string]any{
			"a": 1,
			"b": 2,
			"c": 3,
		},
		"things": []int64{2, 3, 4},
	})
	if err := wtx1.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	wtx2 := rw.Write()
	wtx2.Update(Repository, rr.Key{"id": 2}, rr.Update{
		"value": "World!",
		"list":  []int64{4, 5, 6},
	})
	wtx2.Update(Repository, rr.Key{"id": 3}, rr.Update{
		"value": "Hola?",
	})
	wtx2.Update(Repository, rr.Key{"id": 10}, rr.Update{
		"stuff": rr.Update{
			"b": 4,
		},
	})
	if err := wtx2.Commit(ctx); err != nil {
		t.Fatal(err)
	}

	rtx := rw.Read()
	rtx.Get(Repository, rr.Key{"id": 2})
	rtx.Get(Repository, rr.Key{"id": 3})
	rtx.Get(Repository, rr.Key{"id": 4})
	rtx.Get(Repository, rr.Key{"id": 10})

	res, err := rtx.Execute(ctx)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.MarshalIndent(res, "", "\t")
	t.Log(string(data))
}
