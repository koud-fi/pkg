package transform

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/koud-fi/pkg/file"
)

var (
	paramsParser = regexp.MustCompile("([a-zA-Z]*)([0-9]*)*")
	paramKeys    = []string{"", "x", "t"}
)

type Params map[string]int

func ParseParams(params string) (Params, error) {
	var (
		grps = paramsParser.FindAllStringSubmatch(params, -1)
		p    = make(Params, len(grps))
	)
	for _, g := range grps {
		switch len(g) {
		case 0, 1:
		case 2:
			if v, err := strconv.Atoi(g[1]); err == nil {
				p[""] = v
			}
			p[g[1]] = 0
		default:
			v, _ := strconv.Atoi(g[2])
			p[g[1]] = v
		}
	}
	return p, nil
}

func StdImagePreviewParamsList(attrs file.MediaAttributes) []Params {
	if attrs.Width < 300/0.8 {
		return nil
	}
	var ws []int
	if attrs.Width < 600/0.8 {
		ws = []int{300}
	} else {
		ws = []int{300, 600}
		for {
			l := len(ws)
			if ws[l-1] > int(float64(attrs.Width)*0.8) {
				break
			}
			ws = append(ws, ws[l-2]+ws[l-1])
		}
	}
	ps := make([]Params, 0, len(ws))
	for _, w := range ws {
		ps = append(ps, Params{"": w, "x": 0})
	}
	return ps
}

func (p Params) String() string {
	var sb strings.Builder
	for i, k := range paramKeys {
		sb.WriteString(k)

		v := p[k]
		if v > 0 {
			sb.WriteString(strconv.Itoa(v))
		} else if i < len(paramKeys)-1 {
			sb.WriteString("_")
		}
	}
	return sb.String()
}
