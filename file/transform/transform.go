package transform

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/shell"
)

var (
	defaultImageOnce sync.Once
	defaultImage     []byte
)

type config struct {
	contentType string
	useDefault  bool
}

type Option func(*config)

func ContentType(ct string) Option { return func(c *config) { c.contentType = ct } }
func UseDefault(b bool) Option     { return func(c *config) { c.useDefault = true } }

func ToImage(b blob.Blob, p Params, opt ...Option) blob.Blob {
	return blob.Func(func() (io.ReadCloser, error) {
		var c config
		for _, opt := range opt {
			opt(&c)
		}
		src, contentType, err := srcAndType(b, c.contentType)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		switch strings.SplitN(contentType, "/", 2)[0] {
		case "image":
			return toImage(src, p)
		case "video":
			return videoToImage(src, p)
		default:
			if c.useDefault {
				return generateDefaultImage(), nil
			}
			return nil, fmt.Errorf("unsupported content-type: %s", contentType)
		}
	})
}

func srcAndType(b blob.Blob, contentType string) (path, typ string, err error) {
	var rc io.ReadCloser
	if rc, err = b.Open(); err != nil {
		return
	}
	defer rc.Close()

	if f, ok := rc.(*os.File); ok {
		path, _ = filepath.Abs(f.Name())
	}
	if contentType == "" {
		peek, _ := bufio.NewReaderSize(rc, 512).Peek(512)
		typ = http.DetectContentType(peek)
	}
	return
}

func toImage(src string, p Params) (io.ReadCloser, error) {
	if src == "" {

		// TODO: native implementation

		return nil, errors.New("non-local blobs not supported")
	}
	var (
		args   []interface{}
		w, wOk = p[""]
		h, hOk = p["x"]
		s      string
	)
	if !hOk || (w > 0 && h > 0) {
		args = []interface{}{"--smartcrop", "attention"}
	}
	if wOk && w > 0 {
		s = strconv.Itoa(w)
	}
	if hOk {
		s += "x"
		if h > 0 {
			s += strconv.Itoa(h)
		}
	}
	if s != "" {
		args = append(args, "--size", s)
	}
	return shell.Run(context.TODO(), "vipsthumbnail", append(args, src, "-o", ".jpg")...).Open()
}

func videoToImage(src string, p Params) (io.ReadCloser, error) {
	if src == "" {

		// TODO: native implementation

		return nil, errors.New("non-local blobs not supported")
	}
	args := []interface{}{"-v", "fatal"}
	if seek, _ := p["t"]; seek > 0 {
		var (
			t  = time.Duration(seek) * time.Second
			ts = fmt.Sprintf("%02d:%02d:%02d", int(t.Hours()), int(t.Minutes()), int(t.Seconds()))
		)
		args = append(args, "-ss", ts)
	} else {
		args = append(args, "-ss", "00:00:00")
	}
	args = append(args, "-i", src, "-vframes", "1", "-q:v", "5")
	var (
		w, _ = p[""]
		h, _ = p["x"]
	)
	switch {
	case w > 0 && h > 0:
		args = append(args, "-vf", fmt.Sprintf("crop=%d:%d", w, h))
	case w > 0:
		args = append(args, "-vf", fmt.Sprintf("scale=%d:-1", w))
	case h > 0:
		args = append(args, "-vf", fmt.Sprintf("scale=-1:%d", h))
	}
	return shell.Run(context.TODO(), "ffmpeg", append(args, "-f", "mjpeg", "-")...).Open()
}

func generateDefaultImage() io.ReadCloser {

	// TODO: scale output correctly

	defaultImageOnce.Do(func() {
		buf := bytes.NewBuffer(nil)
		png.Encode(buf, image.NewRGBA(image.Rect(0, 0, 1, 1)))
		defaultImage = buf.Bytes()
	})
	return io.NopCloser(bytes.NewReader(defaultImage))
}
