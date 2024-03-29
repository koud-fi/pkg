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
	"github.com/koud-fi/pkg/file/format/raw"
	"github.com/koud-fi/pkg/shell"
)

var (
	DefaultImageOutputExt = ".jpg"

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
			switch contentType {
			case raw.RAFMime:
				return rafToImage(b, p)
			default:
				return toImage(b, p)
			}
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

func srcAndType(
	b blob.Blob, contentType string,
) (path, typ string, _ error) {
	rc, err := b.Open()
	if err != nil {
		return "", "", err
	}
	if f, ok := rc.(*os.File); ok {
		path, _ = filepath.Abs(f.Name())
	}
	if typ = contentType; typ == "" {
		peek, _ := bufio.NewReaderSize(rc, 512).Peek(512)
		typ = http.DetectContentType(peek)
	}
	return
}

func toImage(b blob.Blob, p Params) (io.ReadCloser, error) {
	var (
		args []any
		s    string
	)
	if p.Width > 0 && p.Height > 0 {
		args = []any{"--smartcrop", "attention"}
	}
	if p.Width > 0 {
		s = strconv.Itoa(p.Width)
	}
	if p.Height >= 0 {
		s += "x"
		if p.Height > 0 {
			s += strconv.Itoa(p.Height)
		}
	}
	if s != "" {
		args = append(args, "--size", s)
	}
	args = append(args, "-o", DefaultImageOutputExt) // output
	args = append(args, "stdin", b)                  // input
	return shell.Run(context.TODO(), "vipsthumbnail", args...).Open()
}

func rafToImage(b blob.Blob, p Params) (io.ReadCloser, error) {
	raf, err := raw.DecodeRAF(b)
	if err != nil {
		return nil, err
	}
	jpegData, err := raf.JPEG(b)
	if err != nil {
		return nil, err
	}
	return toImage(blob.FromBytes(jpegData), p)
}

func videoToImage(src string, p Params) (io.ReadCloser, error) {
	if src == "" {

		// TODO: native implementation

		return nil, errors.New("non-local blobs not supported")
	}
	args := []any{"-hide_banner", "-v", "fatal"}
	if seek := p.AtTimestamp; seek > 0 {
		var (
			t  = time.Duration(seek) * time.Second
			ts = fmt.Sprintf("%02d:%02d:%02d", int(t.Hours()), int(t.Minutes()), int(t.Seconds()))
		)
		args = append(args, "-ss", ts)
	} else {
		args = append(args, "-ss", "00:00:00")
	}
	args = append(args, "-i", src, "-vframes", "1", "-q:v", "5")
	switch {
	case p.Width > 0 && p.Height > 0:
		args = append(args, "-vf", fmt.Sprintf("crop=%d:%d", p.Width, p.Height))
	case p.Width > 0:
		args = append(args, "-vf", fmt.Sprintf("scale=%d:-1", p.Width))
	case p.Height > 0:
		args = append(args, "-vf", fmt.Sprintf("scale=-1:%d", p.Height))
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
