package skipset

import "testing"

func TestOrdered(t *testing.T) {
	testIntSet(t, func() anyskipset[int] {
		return New[int]()
	})
	testIntSetDesc(t, func() anyskipset[int] {
		return NewDesc[int]()
	})
	testStringSet(t, func() anyskipset[string] {
		return New[string]()
	})
}
