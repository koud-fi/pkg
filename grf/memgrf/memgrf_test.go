package memgrf_test

import (
	"testing"

	"github.com/koud-fi/pkg/grf/grftest"
	"github.com/koud-fi/pkg/grf/memgrf"
)

func Test(t *testing.T) {
	grftest.Test(t, memgrf.NewStore())
}
