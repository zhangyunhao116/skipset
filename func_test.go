package skipset

import (
	"testing"
)

func TestFunc(t *testing.T) {
	testIntSet(t, func() anyskipset[int] {
		return NewFunc(func(a, b int) bool {
			return a < b
		})
	})
	testIntSetDesc(t, func() anyskipset[int] {
		return NewFunc(func(a, b int) bool {
			return a > b
		})
	})
	testStringSet(t, func() anyskipset[string] {
		return NewFunc(func(a, b string) bool {
			return a < b
		})
	})
}
