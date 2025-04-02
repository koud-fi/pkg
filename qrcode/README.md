### Usage example
```golang
const size = 256
var (
    payload = "Hello, world!"
    bgImg   = image.NewRGBA(image.Rect(0, 0, size, size))
    drawFn  = qrcode.DrawSquarePx(color.Black, qrcode.UniformPx(1))
)
draw.Draw(bgImg, bgImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

img, err := qrcode.Draw(payload, size, bgImg, drawFn)
// ...
```
