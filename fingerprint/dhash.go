package fingerprint

import (
	"fmt"
	"image"
	"math/bits"

	"golang.org/x/image/draw"
)

type DHash uint64

func NewDHash(img image.Image) DHash { return newDHash(img, draw.CatmullRom) }

func newDHash(img image.Image, scaleWith *draw.Kernel) DHash {
	out := image.NewGray(image.Rect(0, 0, 8, 9))
	scaleWith.Scale(out, out.Bounds(), img, img.Bounds(), draw.Src, nil)

	var h uint64
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if out.GrayAt(x, y).Y > out.GrayAt(x, y+1).Y {
				h |= 1 << uint(y*8+x)
			}
		}
	}
	return DHash(h)
}

func (h DHash) Distance(to DHash) float64 {
	return float64(bits.OnesCount64(uint64(h)^uint64(to))) / 64.0
}

func (h DHash) String() string { return fmt.Sprintf("%016x", uint64(h)) }
