package endpoint

import (
	"fmt"
	"io"
	"net/http"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/serve"
)

type redirect struct {
	URL    string
	Status int
}

func Redirect(url string, status int) redirect {
	return redirect{URL: url, Status: status}
}

func serveOutput(
	w http.ResponseWriter, r *http.Request, output any, opts []serve.Option,
) (*serve.Info, error) {
	switch v := output.(type) {
	case nil:
		return serve.Blob(w, r, blob.Empty(), opts...)
	case io.ReadCloser:
		defer v.Close()
		return serve.Reader(w, r, v, opts...)
	case io.Reader:
		return serve.Reader(w, r, v, opts...)
	case blob.Blob:
		return serve.Blob(w, r, v, opts...)
	case []byte:
		return serve.Blob(w, r, blob.FromBytes(v), opts...)
	case string:
		return serve.Blob(w, r, blob.FromString(v), opts...)
	case fmt.Stringer:
		return serve.Blob(w, r, blob.FromString(v.String()), opts...)
	case redirect:
		return serve.Redirect(w, r, v.URL, v.Status)
	default:
		return serve.JSON(w, r, v, opts...)
	}
}
