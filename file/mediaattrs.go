package file

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

// TODO: resolution
