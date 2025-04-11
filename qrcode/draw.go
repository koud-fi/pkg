package qrcode

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/koud-fi/pkg/noise/simplex"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/skip2/go-qrcode"
	"golang.org/x/image/draw"
)

// TODO: base all drawing logic to gfx package once it is complete enough

func Draw(payload string, size int, bgImg image.Image, drawFn DrawPxFunc) (image.Image, error) {
	if bgImg == nil {
		img := image.NewRGBA(image.Rect(0, 0, size, size))

		// TODO: make this configurable (drawFn should include some "foreground" color info?)

		draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
		bgImg = img
	}
	mask, err := NewQRMask(payload, size, drawFn)
	if err != nil {
		return nil, err
	}
	dst := image.NewRGBA(mask.Bounds())
	draw.CatmullRom.Scale(dst, dst.Bounds(), bgImg, bgImg.Bounds(), draw.Src, nil)

	for x := range dst.Rect.Dx() {
		for y := range dst.Rect.Dy() {
			var (
				dr, dg, db, da = dst.At(x, y).RGBA()
				mr, mg, mb, ma = mask.At(x, y).RGBA()
			)
			dst.SetRGBA64(x, y, color.RGBA64{
				uint16(dr * mr / 0xffff),
				uint16(dg * mg / 0xffff),
				uint16(db * mb / 0xffff),
				uint16(da * ma / 0xffff),
			})
		}
	}
	return dst, nil
}

type DrawPxFunc func(gc *draw2dimg.GraphicContext, getPx func(x, y int) bool, x, y int, pxSize float64)

func DrawSquarePx(col color.Color, sizeFn PxSizeFn) DrawPxFunc {
	return func(gc *draw2dimg.GraphicContext, getPx func(x, y int) bool, x, y int, pxSize float64) {
		if getPx(x, y) {
			wh := sizeFn(x, y) * pxSize
			DrawSquare(gc, col, (float64(x)+0.5)*pxSize, (float64(y)+0.5)*pxSize, wh, wh)
		}
	}
}

func DrawSquare(gc *draw2dimg.GraphicContext, col color.Color, x, y, w, h float64) {
	gc.SetFillColor(col)
	gc.BeginPath()
	w /= 2
	h /= 2
	gc.MoveTo(x-w, y-h)
	gc.LineTo(x+w, y-h)
	gc.LineTo(x+w, y+h)
	gc.LineTo(x-w, y+h)
	gc.Close()
	gc.FillStroke()
}

func DrawBallPx(col color.Color, sizeFn PxSizeFn) DrawPxFunc {
	return func(gc *draw2dimg.GraphicContext, getPx func(x, y int) bool, x, y int, pxSize float64) {
		if getPx(x, y) {
			DrawBall(gc, col, (float64(x)+0.5)*pxSize, (float64(y)+0.5)*pxSize, pxSize*sizeFn(x, y)*0.5)
		}
	}
}

func DrawBall(gc *draw2dimg.GraphicContext, col color.Color, x, y, size float64) {
	gc.SetFillColor(col)
	gc.BeginPath()
	gc.ArcTo(x, y, size, size, 0, 2*math.Pi)
	gc.Close()
	gc.FillStroke()
}

type PxSizeFn func(x, y int) float64

func UniformPx(size float64) PxSizeFn {
	return func(_, _ int) float64 { return size }
}

func RandPx(seed int64, minSize, maxSize float64) PxSizeFn {
	rng := rand.New(rand.NewSource(seed))
	return func(_, _ int) float64 {

		// TODO: factor location into the randomness instead of just using the next value

		return minSize + rng.Float64()*(maxSize-minSize)
	}
}

func NoisePx(seed int64, minSize, maxSize float64) PxSizeFn {
	return func(x, y int) float64 {

		// TODO: improve the noise function, especially the seed handling

		n := float64(simplex.Noise2D(float32(int64(x)+seed)*32, float32(int64(y)+seed*-1)*32)/2 + 0.5)
		return minSize + n*(maxSize-minSize)
	}
}

func NewQRMask(data string, size int, drawFn DrawPxFunc) (image.Image, error) {
	qrc, err := qrcode.New(data, qrcode.Highest)
	if err != nil {
		return nil, fmt.Errorf("could not generate QRCode: %w", err)
	}
	qrc.DisableBorder = false

	var (
		bm = qrc.Bitmap()
		//borderWidth = 4 // TODO: resolve correct value based on the QR code version
		mask          = image.NewRGBA(image.Rect(0, 0, size, size))
		gc            = draw2dimg.NewGraphicContext(mask)
		bgLum   uint8 = 255
		bgColor       = color.RGBA{bgLum, bgLum, bgLum, 255}
	)
	draw.Draw(mask, mask.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)
	gc.SetStrokeColor(color.RGBA{0, 0, 0, 255})
	gc.SetLineWidth(0)

	pxSize := float64(size) / float64(len(bm))
	for i := 0; i < len(bm); i++ {
		row := bm[i]
		for j := 0; j < len(row); j++ {
			drawFn(gc, func(x, y int) bool {

				// TODO: bounds checking

				return bm[j][i]
			}, j, i, pxSize)
		}
	}
	return mask, nil
}
