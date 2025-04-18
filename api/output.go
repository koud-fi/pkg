package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/serve"
)

type JSONWrapper[T any] struct {
	Ok    bool  `json:"ok"`
	Error error `json:"error,omitempty"`
	Data  T     `json:"data,omitempty"`
}

func serveOutput(
	w http.ResponseWriter, r *http.Request, output any, err error,
) (*serve.Info, error) {
	switch v := output.(type) {
	case nil:
		if err != nil {
			return nil, err
		}
		return serve.Blob(w, r, blob.Empty())
	case io.ReadCloser:
		defer v.Close()
		return serve.Reader(w, r, v)
	case io.Reader:
		return serve.Reader(w, r, v)
	case blob.Blob:
		return serve.Blob(w, r, v)
	case []byte:
		return serve.Blob(w, r, blob.FromBytes(v))
	case string:
		return serve.Blob(w, r, blob.FromString(v))
	case fmt.Stringer:
		return serve.Blob(w, r, blob.FromString(v.String()))
	default:
		return serve.JSON(w, r, JSONWrapper[any]{
			Ok:    true,
			Data:  v,
			Error: err,
		})
	}
}
