package endpoint

import (
	"fmt"
	"io"
	"net/http"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/serve"
)

type JSONWrapper[T any] struct {
	Ok    bool  `json:"ok"`
	Error error `json:"error"`
	Data  T     `json:"data"`
}

func serveOutput(
	w http.ResponseWriter, r *http.Request, output any, err error,
) (*serve.Info, error) {
	if err != nil {
		return serve.JSON(w, r,
			JSONWrapper[any]{Ok: false, Error: err},
			serve.StatusCode(http.StatusBadRequest),
		)
	}
	switch v := output.(type) {
	case nil:
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
		return serve.JSON(w, r, JSONWrapper[any]{Ok: true, Data: v})
	}
}
