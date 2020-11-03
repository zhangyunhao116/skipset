package skipset

import "github.com/zeebo/xxh3"

// TODO
func xxh3Hash(s string) uint64 {
	return xxh3.HashString(s)
}
