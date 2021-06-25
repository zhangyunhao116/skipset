package skipset

import (
	_ "unsafe" // for linkname
)

const (
	maxLevel            = 16
	p                   = 0.25
	defaultHighestLevel = 3
)

//go:linkname fastrand runtime.fastrand
func fastrand() uint32

//go:linkname cmpstring runtime.cmpstring
func cmpstring(a, b string) int

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
