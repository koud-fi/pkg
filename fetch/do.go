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

func (r *Request) Open() (io.ReadCloser, error) { return r.OpenFile() }

func (r *Request) OpenFile() (fs.File, error) {
	res, req, err := r.do()
	if err != nil {
		return nil, err
	}
	return &file{
		fileInfo: &fileInfo{url: req.URL, header: res.Header, dr: r.dirReader},
		body:     res.Body,
	}, nil
}

func (r *Request) Stat() (fs.FileInfo, error) {
	res, req, err := r.do()
	if err != nil {
		return nil, err
	}
	res.Body.Close()
	return &fileInfo{url: req.URL, header: res.Header, dr: r.dirReader}, nil
}

func (r *Request) do() (*http.Response, *http.Request, error) {
	req, err := r.HttpRequest()
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
	if res.StatusCode >= 400 {
		return nil, nil, processErrorResponse(req, res)
	}
	return res, req, nil
}

func processErrorResponse(req *http.Request, res *http.Response) error {
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
