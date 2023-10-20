package transform

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/koud-fi/pkg/file"
)

var paramsParser = regexp.MustCompile("([a-zA-Z]*)([0-9]*)*")

type Params struct {
	Width       int
	Height      int
	AtTimestamp int
}

func ParseParams(params string) (Params, error) {
	var p Params
	err := processParams(params, func(key, value string) (err error) {
		switch key {
		case "":
			p.Width, err = strconv.Atoi(value)
		case "x":
			p.Height, err = strconv.Atoi(value)
		case "t":
			p.AtTimestamp, err = strconv.Atoi(value)
		}
		return
	})
	return p, err
}

func processParams(params string, process func(key, value string) error) error {
	grps := paramsParser.FindAllStringSubmatch(params, -1)
	for _, g := range grps {
		switch len(g) {
		case 0, 1:
		case 2:
			process("", g[1])
		default:
			process(g[1], g[2])
		}
	}
	return nil
}

func (p Params) String() string {
	var sb strings.Builder
	if p.Width > 0 {
		sb.WriteString(strconv.Itoa(p.Width))
	}
	if p.Height > 0 {
		sb.WriteString("x" + strconv.Itoa(p.Height))
	}
	if p.AtTimestamp > 0 {
		sb.WriteString("t" + strconv.Itoa(p.AtTimestamp))
	}
	return sb.String()
}

func StdImagePreviewParams(attrs file.MediaAttributes) map[Params]file.MediaAttributes {
	if attrs.Width < 300/0.8 {
		return nil
	}
	var ws []int
	if attrs.Width < 600/0.8 {
		ws = []int{300}
	} else {
		ws = []int{300, 600}
		for {
			var (
				l = len(ws)
				w = ws[l-2] + ws[l-1]
			)
			if w > int(float64(attrs.Width)*0.8) {
				break
			}
			ws = append(ws, w)
		}
	}
	ps := make([]Params, 0, len(ws))
	for _, w := range ws {
		ps = append(ps, Params{Width: w})
	}
	out := make(map[Params]file.MediaAttributes, len(ps))
	for _, p := range ps {
		mod := float64(p.Width) / float64(attrs.Width)
		out[p] = file.MediaAttributes{
			Width:  p.Width,
			Height: int(float64(attrs.Height) * mod),
		}
	}
	return out
}
