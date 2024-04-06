package dice

import (
	"fmt"
	"strings"

	"github.com/koud-fi/pkg/num"
)

var unicodeDie = []rune{'⚀', '⚁', '⚂', '⚃', '⚄', '⚅'}

type Result struct {
	Die   Die
	Rolls []int
}

func (r Result) Total() int { return num.Sum(r.Rolls...) }

func (r Result) String() string {
	var (
		sb  strings.Builder
		sum int
	)
	for i, roll := range r.Rolls {
		if r.Die == D6 {
			if i > 0 {
				sb.WriteByte(',')
			}
			sb.WriteRune(unicodeDie[roll-1])
		} else {
			if i > 0 {
				sb.WriteByte('+')
			}
			fmt.Fprint(&sb, roll)
		}
		sum += roll
	}
	fmt.Fprintf(&sb, " = %d", sum)
	return sb.String()
}
