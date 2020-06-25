package skipset

import (
	"math/rand"
)

const (
	maxLevel = 32
	p        = 0.25
)

func randomLevel() int {
	level := 1
	for rand.Float64() < p {
		level++
	}
	if level > maxLevel {
		return maxLevel
	}
	return level
}
