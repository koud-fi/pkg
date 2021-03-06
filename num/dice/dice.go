package dice

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/koud-fi/pkg/num"
)

var D20 Die = 20

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

type Dice struct {
	N   int
	Die Die
}

func Parse(s string) (d Dice) {
	d.UnmarshalJSON([]byte(s))
	return
}

func (t Dice) Roll(mod, advantage int) []int {
	rolls := make([]int, t.N)
	for i := range rolls {
		rolls[i] = t.Die.Roll(mod, advantage)
	}
	return rolls
}

func (t Dice) Average(mod int) int {
	return int(math.Round(float64(t.Max(mod)+t.N) / 2))
}

func (t Dice) Max(mod int) int {
	sum := 0
	for i := 0; i < t.N; i++ {
		sum += t.Die.Max(mod)
	}
	return sum
}

func (t Dice) String() string { return fmt.Sprintf(`"%dd%d"`, t.N, t.Die) }

func (t Dice) MarshalJSON() ([]byte, error) { return []byte(t.String()), nil }

func (t *Dice) UnmarshalJSON(data []byte) error {
	parts := strings.Split(strings.Trim(string(data), `"`), "d")
	if len(parts) == 1 {
		die, _ := strconv.Atoi(parts[0])
		t.Die = Die(die)
	} else {
		t.N, _ = strconv.Atoi(parts[0])
		die, _ := strconv.Atoi(parts[1])
		t.Die = Die(die)
	}
	if t.N <= 0 {
		t.N = 1
	}
	return nil
}
