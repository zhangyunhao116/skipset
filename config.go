package skipset

import (
	"math"
	"math/rand"

	"github.com/ZYunH/lockedsource"
)

const (
	maxLevel = 32
	p        = 1 / math.E
)

var rnd = rand.New(lockedsource.New(0))

func randomLevel() int {
	level := 1
	for rnd.Float64() < p {
		level++
	}
	if level > maxLevel {
		return maxLevel
	}
	return level
}
