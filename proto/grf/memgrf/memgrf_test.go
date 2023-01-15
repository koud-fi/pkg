package memgrf_test

import (
	"testing"

	"github.com/koud-fi/pkg/proto/grf/grftest"
	"github.com/koud-fi/pkg/proto/grf/memgrf"
)

func Test(t *testing.T) {
	grftest.Test(t, memgrf.NewStore())
}
