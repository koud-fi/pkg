package dice

import (
	"math/rand"

	"github.com/koud-fi/pkg/num"
)

const (
	D4   Die = 4
	D6   Die = 6
	D8   Die = 8
	D10  Die = 10
	D12  Die = 12
	D20  Die = 20
	D100 Die = 100
)

type Die int

func (t Die) Roll(mod, advantage int) int {
	switch {
	case advantage < 0:
		return num.Min((Dice{1 + -advantage, t}.Roll(mod, 0))...)
	case advantage > 0:
		return num.Max((Dice{1 + advantage, t}.Roll(mod, 0))...)
	default:
		return t.roll(mod)
	}
}

func (t Die) roll(mod int) int {
	if t <= 0 {
		return t.Max(mod)
	}
	return 1 + rand.Intn(int(t)) + mod
}

func (t Die) Max(mod int) int { return int(t) + mod }
