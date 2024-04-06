package dice

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var alias = map[string]Dice{
	"fireball": {N: 8, Die: D6},
}

type Dice struct {
	N   int
	Die Die
}

func Parse(s string) (d Dice, err error) {
	err = d.UnmarshalJSON([]byte(s))
	return
}

func (d Dice) Roll(mod, advantage int) Result {
	rolls := make([]int, d.N)
	for i := range rolls {
		rolls[i] = d.Die.Roll(mod, advantage)
	}
	return Result{
		Die:   d.Die,
		Rolls: rolls,
	}
}

func (d Dice) Average(mod int) int {
	return int(math.Round(float64(d.Max(mod)+d.N) / 2))
}

func (d Dice) Max(mod int) int {
	sum := 0
	for i := 0; i < d.N; i++ {
		sum += d.Die.Max(mod)
	}
	return sum
}

func (d Dice) String() string { return fmt.Sprintf(`"%dd%d"`, d.N, d.Die) }

func (d Dice) MarshalJSON() ([]byte, error) { return []byte(d.String()), nil }

func (d *Dice) UnmarshalJSON(data []byte) error {
	s := string(data)

	if dice, ok := alias[strings.ToLower(s)]; ok {
		*d = dice
		return nil
	}
	parts := strings.Split(strings.Trim(s, `"`), "d")
	if len(parts) == 1 {
		die, _ := strconv.Atoi(parts[0])
		d.Die = Die(die)
	} else {
		d.N, _ = strconv.Atoi(parts[0])
		die, _ := strconv.Atoi(parts[1])
		d.Die = Die(die)
	}
	if d.N <= 0 {
		d.N = 1
	}
	return nil
}
