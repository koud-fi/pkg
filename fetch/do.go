package fetch

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
)

const maxErrLen = 1 << 12

func (r *Request) openFile() (fs.File, error) {
	res, req, err := r.do()
	if err != nil {
		return nil, err
	}
	return &file{
		fileInfo: newFileInfo(r, res, req),
		body:     res.Body,
	}, nil
}

func (r *Request) stat() (fs.FileInfo, error) {
	res, req, err := r.do()
	if err != nil {
		return nil, err
	}
	res.Body.Close()
	return newFileInfo(r, res, req), nil
}

func (r *Request) do() (*http.Response, *http.Request, error) {
	req, err := r.httpRequest()
	if err != nil {
		return nil, nil, err
	}
	if r.limiter != nil {
		if err := r.limiter.Wait(context.Background()); err != nil {
			return nil, nil, err
		}
	}
	res, err := r.client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	// TODO: don't automatically treat all 400 >= statuses as errors
	// TODO: add custom error decoders?

	if res.StatusCode >= 400 {
		return nil, nil, processErrorResponse(req, res)
	}
	return res, req, nil
}

func processErrorResponse(_ *http.Request, res *http.Response) error {
	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read error body: %w", err)
	}

	// TODO: handle binary response bodies gracefully

	var msg string
	if len(buf) > maxErrLen {
		msg = string(buf[:maxErrLen]) + "..."
	} else {
		msg = string(buf)
	}
	switch res.StatusCode {
	case http.StatusNotFound:
		return fmt.Errorf("%w: %s", fs.ErrNotExist, msg)
	case http.StatusForbidden:
		return fmt.Errorf("%w: %s", fs.ErrPermission, msg)
	case http.StatusGatewayTimeout:
		return fmt.Errorf("%w: %s", os.ErrDeadlineExceeded, msg)
	default:
		return fmt.Errorf("%s: %s", res.Status, msg)
	}
}
