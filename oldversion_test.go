package skipset

import (
	"testing"
)

func TestOldVersion(t *testing.T) {
	testIntSet(t, func() anyskipset[int] {
		return NewInt()
	})
	testIntSetDesc(t, func() anyskipset[int] {
		return NewIntDesc()
	})
	// testStringSet(t, func() anyskipset[string] {
	// 	return NewString()
	// })
}
