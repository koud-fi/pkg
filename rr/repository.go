package rr

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not round")

type Repository string
type Item map[string]any
type Key = Item
type Update = Item

type ReadWriter interface {
	Reader
	Writer
}

type Reader interface {
	Read() ReadTx
}

type ReadTx interface {
	Get(Repository, Key)
	Execute(ctx context.Context) (map[Repository][]Item, error)
}

type Writer interface {
	Write() WriteTx
}

type WriteTx interface {
	Put(Repository, Item)
	Update(Repository, Key, Update)
	Delete(Repository, Key)
	Commit(context.Context) error
}
