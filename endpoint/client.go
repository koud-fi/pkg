package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/koud-fi/pkg/blob"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func (c *Client) Invoke(ctx context.Context, endpoint string, in any) blob.Reader {
	return blob.Func(func() (io.ReadCloser, error) {

		// TODO: Use the "modern" fetch package once that is ported over (ask Pasi for details)

		body, err := json.Marshal(in)
		if err != nil {
			return nil, fmt.Errorf("marshal input: %w", err)
		}
		var (
			url        = c.baseURL + endpoint
			bodyReader = bytes.NewReader(body)
		)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
		if err != nil {
			return nil, fmt.Errorf("new request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		res, err := c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("do request: %w", err)
		}
		return res.Body, nil
	})
}
