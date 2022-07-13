package simplex_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/koud-fi/pkg/noise/simplex"
)

func TestNoise1D(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 1024, 32))
	for x := 0; x < img.Rect.Dx(); x++ {
		n := simplex.Noise1D(float32(x) / float32(img.Rect.Dx()/8))
		for y := 0; y < img.Rect.Dy(); y++ {
			img.SetGray(x, y, color.Gray{Y: uint8((n + 1) * 127)})
		}
	}
	writeImageFile(t, "temp/simplex_1d.png", img)
}

func TestNoise2D(t *testing.T) {
	var (
		img = image.NewGray(image.Rect(0, 0, 1024, 1024))
		fdx = float32(img.Rect.Dx())
		fdy = float32(img.Rect.Dy())
	)
	for x := 0; x < img.Rect.Dx(); x++ {
		fx := float32(x) / fdx
		for y := 0; y < img.Rect.Dy(); y++ {
			fy := float32(y) / fdy
			n := simplex.Noise2D(fx, fy) +
				0.5*simplex.Noise2D(2*fx, 2*fy) +
				0.25*simplex.Noise2D(4*fx, 4*fy) +
				0.1*simplex.Noise2D(8*fx, 8*fy) +
				0.025*simplex.Noise2D(16*fx, 16*fy)
			n /= 1.875
			img.SetGray(x, y, color.Gray{Y: uint8((n + 1) * 127)})
		}
	}
	writeImageFile(t, "temp/simplex_2d.png", img)
}

func writeImageFile(t *testing.T, path string, img image.Image) {
	buf := bytes.NewBuffer(nil)
	if err := png.Encode(buf, img); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
}
