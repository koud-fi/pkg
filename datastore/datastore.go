package datastore

import "github.com/koud-fi/pkg/blob"

type Store struct {
	s blob.Storage
	c Codec
}

type Codec struct {
	Marshal   blob.MarshalFunc
	Unmarshal blob.UnmarshalFunc
}

func New(s blob.Storage, c Codec) *Store {
	return &Store{s: s}
}

// TODO
