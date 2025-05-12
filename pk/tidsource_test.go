package pk_test

import (
	"testing"
	"time"

	"github.com/koud-fi/pkg/pk"
)

func TestTIDSource(t *testing.T) {

	// TODO: Improve, run multiple sources in parallel and check for collisions.

	src := pk.NewTIDSource(time.Now(), 1)
	for range 1 << 13 {
		src.Next(time.Now())
	}
}
