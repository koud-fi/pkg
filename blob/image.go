package blob

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"sync"
)

type (
	ImageFunc    func() (image.Image, error)
	ImageEncoder func(w io.Writer, img image.Image) error
)

func (fn ImageFunc) Reader(encode ImageEncoder) Reader {
	return &imageReader{imgFn: fn, encode: encode}
}

func (fn ImageFunc) PNG() Reader {
	return fn.Reader(png.Encode)
}

func (fn ImageFunc) JPEG(quality int) Reader {
	return fn.Reader(func(w io.Writer, img image.Image) error {
		return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
	})
}

type imageReader struct {
	imgFn  ImageFunc
	encode ImageEncoder

	once sync.Once
	buf  bytes.Buffer
	err  error
}

func (ir *imageReader) Open() (io.ReadCloser, error) {
	ir.once.Do(func() {
		img, err := ir.imgFn()
		if err != nil {
			ir.err = err
			return
		}
		ir.err = ir.encode(&ir.buf, img)
	})
	if ir.err != nil {
		return nil, ir.err
	}
	return io.NopCloser(bytes.NewReader(ir.buf.Bytes())), nil
}
