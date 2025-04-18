package endpoint

import (
	"errors"
	"fmt"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/fetch"
)

func Fetch[T any](r *fetch.Request, in any) (T, error) {
	var zero T
	out, err := blob.UnmarshalJSONValue[JSONWrapper[T]](r.JSON(in))
	if err != nil {
		return zero, fmt.Errorf("fetch: %w", err)
	}
	if !out.Ok {
		if out.Error == nil {
			return zero, errors.New("fail without error")
		}
		return zero, out.Error
	}
	return out.Data, nil
}
