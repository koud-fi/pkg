package dice

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

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
