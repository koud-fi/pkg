package blob

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"sync"
)

type ImageEncoder func(w io.Writer, img image.Image) error

func FromImage(img image.Image, encode ImageEncoder) Reader {
	return &imageReader{img: img, encode: encode}
}

func FromImagePNG(img image.Image) Reader {
	return FromImage(img, png.Encode)
}

func FromImageJPEG(img image.Image, quality int) Reader {
	return FromImage(img, func(w io.Writer, img image.Image) error {
		return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
	})
}

type imageReader struct {
	img    image.Image
	encode ImageEncoder

	once sync.Once
	buf  bytes.Buffer
	err  error
}

func (ir *imageReader) Open() (io.ReadCloser, error) {
	ir.once.Do(func() {
		ir.err = ir.encode(&ir.buf, ir.img)
	})
	if ir.err != nil {
		return nil, ir.err
	}
	return io.NopCloser(bytes.NewReader(ir.buf.Bytes())), nil
}
