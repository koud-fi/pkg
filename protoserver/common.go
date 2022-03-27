package protoserver

import (
	"context"

	"github.com/koud-fi/pkg/pk"
)

const EchoScheme = "echo"

func RegisterEcho() {
	Register(EchoScheme, FetchFunc(func(_ context.Context, ref pk.Ref) (any, error) {
		return ref.Key(), nil
	}))
}
