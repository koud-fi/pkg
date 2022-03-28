package protoserver

import (
	"context"
	"fmt"

	"github.com/koud-fi/pkg/fetch"
	"github.com/koud-fi/pkg/pk"
)

const (
	HTTPScheme  = "http"
	HTTPSScheme = "https"
)

func RegisterHTTP()  { Register(HTTPScheme, FetchFunc(fetchHTTP)) }
func RegisterHTTPS() { Register(HTTPSScheme, FetchFunc(fetchHTTP)) }

func fetchHTTP(ctx context.Context, ref pk.Ref) (any, error) {
	url := fmt.Sprintf("%s://%s", ref.Scheme(), ref.Key())
	return fetch.Get(url).Context(ctx).Open()
}
