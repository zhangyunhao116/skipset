package skipset

import (
	_ "unsafe" // for runtime.fastrand
)

const (
	maxLevel = 32
	p        = 0.25
)

//go:linkname fastrand runtime.fastrand
func fastrand() uint32

//go:nosplit
func fastrandn(n uint32) uint32 {
	// This is similar to fastrand() % n, but faster.
	// See https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/
	return uint32(uint64(fastrand()) * uint64(n) >> 32)
}

func randomLevel() int {
	level := 1
	for fastrandn(1/p) == 0 {
		level++
	}
	if level > maxLevel {
		return maxLevel
	}
	return level
}
