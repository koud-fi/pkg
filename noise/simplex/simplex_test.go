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
		n := simplex.Noise1D(float32(x) / float32(img.Rect.Dx()/4))
		for y := 0; y < img.Rect.Dy(); y++ {
			img.SetGray(x, y, color.Gray{Y: uint8(n * 255)})
		}
	}
	writeImageFile(t, "temp/simplex_1d.png", img)
}

func TestNoise2D(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 1024, 1024))
	for x := 0; x < img.Rect.Dx(); x++ {
		for y := 0; y < img.Rect.Dy(); y++ {
			n := simplex.Noise2D(
				float32(x)/float32(img.Rect.Dx()/4),
				float32(y)/float32(img.Rect.Dy()/4))
			img.SetGray(x, y, color.Gray{Y: uint8(n * 255)})
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
