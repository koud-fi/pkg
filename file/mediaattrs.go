package file

import (
	"context"
	"encoding/json"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/koud-fi/pkg/blob"
	"github.com/koud-fi/pkg/file/format/raw"
	"github.com/koud-fi/pkg/shell"
)

type MediaAttributes struct {
	Width      int     `json:"width,omitempty"`
	Height     int     `json:"height,omitempty"`
	Duration   float64 `json:"duration,omitempty"`
	HasAudio   bool    `json:"hasAudio,omitempty"`
	FrameCount int     `json:"frameCount,omitempty"`
}

func (ma MediaAttributes) AspectRatio() float64 {
	if ma.Width <= 0 || ma.Height <= 0 {
		return 1.0
	}
	return float64(ma.Width) / float64(ma.Height)
}

func (ma MediaAttributes) Megapixels() float64 {
	return float64(ma.Width) * float64(ma.Height) / 1_000_000
}

func MediaAttrs() Option {
	return func(a *Attributes, b blob.Blob, contentType string) error {
		switch contentType {
		case "image/jpeg", "image/png", "image/webp", "image/bmp", "image/tiff":
			return resolveImageAttrs(&a.MediaAttributes, b)
		case "image/gif":
			return resolveGIFAttrs(&a.MediaAttributes, b)
		case "video/mp4", "video/webm":
			return resolveVideoAttrs(&a.MediaAttributes, b)

		// TODO: common audio formats

		case raw.RAFMime:
			return resolveRAFAttrs(&a.MediaAttributes, b)
		}
		return nil
	}
}

func resolveImageAttrs(a *MediaAttributes, b blob.Blob) error {
	return blob.Use(b, func(r io.Reader) error {
		c, _, err := image.DecodeConfig(r)
		if err != nil {
			switch err.(type) {
			case jpeg.FormatError, png.FormatError:
				return nil
			}
			return err
		}
		a.Width = c.Width
		a.Height = c.Height
		return nil
	})
}

func resolveGIFAttrs(a *MediaAttributes, b blob.Blob) error {
	return blob.Use(b, func(r io.Reader) error {
		g, err := gif.DecodeAll(r)
		if err != nil {
			return err
		}
		if len(g.Image) == 0 {
			return nil
		}
		for _, delay := range g.Delay {
			a.Duration += float64(delay) / 100.0
		}
		a.Width = g.Image[0].Bounds().Dx()
		a.Height = g.Image[0].Bounds().Dy()
		a.FrameCount = len(g.Image)
		return nil
	})
}

func resolveRAFAttrs(a *MediaAttributes, b blob.Blob) error {

	// TODO

	return nil
}

func resolveVideoAttrs(a *MediaAttributes, b blob.Blob) error {

	// TODO: native implementation

	var info ffprobeInfo
	if err := blob.Unmarshal(json.Unmarshal, shell.Run(context.TODO(), "ffprobe",
		"-i", "-", b,
		"-v", "fatal",
		"-of", "json",
		"-show_format", "-show_streams",
	), &info); err != nil {
		return err
	}
	vstream := info.findStream("video")
	if vstream == nil {
		return nil
	}
	a.Width = vstream.Width
	a.Height = vstream.Height
	a.Duration = info.Format.Duration
	a.HasAudio = info.findStream("audio") != nil
	return nil
}

type ffprobeInfo struct {
	Streams []ffprobeSteam
	Format  ffprobeFormat
}

func (t ffprobeInfo) findStream(codecType string) *ffprobeSteam {
	for _, s := range t.Streams {
		if s.CodecType == codecType {
			return &s
		}
	}
	return nil
}

type ffprobeSteam struct {
	CodecType string `json:"codec_type"`
	Width     int
	Height    int
}

type ffprobeFormat struct {
	Duration float64 `json:",string"`
}
