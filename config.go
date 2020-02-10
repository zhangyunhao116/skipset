package skipset

import (
	"math/rand"

	"github.com/ZYunH/lockedsource"
)

const (
	maxLevel = 32
	p        = 0.25
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
