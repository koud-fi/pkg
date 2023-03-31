package mongo_test

import (
	"testing"

	"github.com/koud-fi/pkg/rr/mongo"
	"github.com/koud-fi/pkg/rr/rrtest"
)

func TestRW(t *testing.T) {

	// TODO

	rrtest.Run(t, mongo.NewRW(nil))
}
