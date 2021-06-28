package skipset

import (
	_ "unsafe" // for linkname

	"github.com/zhangyunhao116/fastrand"
	"github.com/zhangyunhao116/wyhash"
)

const (
	maxLevel            = 16
	p                   = 0.25
	defaultHighestLevel = 3
)

//go:linkname cmpstring runtime.cmpstring
func cmpstring(a, b string) int

func hash(s string) uint64 {
	return wyhash.Sum64String(s)
}

func randomLevel() int {
	level := 1
	for fastrand.Uint32n(1/p) == 0 {
		level++
	}
	if level > maxLevel {
		return maxLevel
	}
	return level
}
