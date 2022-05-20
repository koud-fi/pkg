package datastore

import "github.com/koud-fi/pkg/blob"

type Store struct {
	s blob.Storage
}

func New(s blob.Storage) *Store {
	return &Store{s: s}
}

// TODO
