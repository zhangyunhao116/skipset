// Package skipset is a high-performance, scalable, concurrent-safe set based on skip-list.
// In the typical pattern(100000 operations, 90%CONTAINS 9%Add 1%Remove, 8C16T), the skipset
// up to 15x faster than the built-in sync.Map.
//
//go:generate go run gen.go
package skipset

import "math"

// New returns an empty skip set in ascending order.
func New[T ordered]() *OrderedSet[T] {
	var t T
	h := newOrderedNode(t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &OrderedSet[T]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewDesc returns an empty skip set in descending order.
func NewDesc[T ordered]() *OrderedSetDesc[T] {
	var t T
	h := newOrderedNodeDesc(t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &OrderedSetDesc[T]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewFunc returns an empty skip set in ascending order.
//
// Note that the less function requires a strict weak ordering,
// see https://en.wikipedia.org/wiki/Weak_ordering#Strict_weak_orderings,
// or undefined behavior will happen.
func NewFunc[T any](less func(a, b T) bool) *FuncSet[T] {
	var t T
	h := newFuncNode(t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &FuncSet[T]{
		header:       h,
		highestLevel: defaultHighestLevel,
		less:         less,
	}
}

// NewString returns an empty skip set in ascending order.
func NewString() *StringSet {
	h := newStringNode("", maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &StringSet{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewStringDesc returns an empty skip set in descending order.
func NewStringDesc() *StringSetDesc {
	h := newStringNodeDesc("", maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &StringSetDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

func isNaNf32(x float32) bool {
	return x != x
}

// NewFloat32 returns an empty skip set in ascending order.
func NewFloat32() *FuncSet[float32] {
	return NewFunc(func(a, b float32) bool {
		return a < b || (isNaNf32(a) && !isNaNf32(b))
	})
}

// NewFloat32Desc returns an empty skip set in descending order.
func NewFloat32Desc() *FuncSet[float32] {
	return NewFunc(func(a, b float32) bool {
		return a > b || (isNaNf32(a) && !isNaNf32(b))
	})
}

// NewFloat64 returns an empty skip set in ascending order.
func NewFloat64() *FuncSet[float64] {
	return NewFunc(func(a, b float64) bool {
		return a < b || (math.IsNaN(a) && !math.IsNaN(b))
	})
}

// NewFloat64Desc returns an empty skip set in descending order.
func NewFloat64Desc() *FuncSet[float64] {
	return NewFunc(func(a, b float64) bool {
		return a > b || (math.IsNaN(a) && !math.IsNaN(b))
	})
}

// NewInt returns an empty skip set in ascending order.
func NewInt() *IntSet {
	h := newIntNode(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &IntSet{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewIntDesc returns an empty skip set in descending order.
func NewIntDesc() *IntSetDesc {
	h := newIntNodeDesc(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &IntSetDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewInt64 returns an empty skip set in ascending order.
func NewInt64() *Int64Set {
	h := newInt64Node(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Int64Set{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewInt64Desc returns an empty skip set in descending order.
func NewInt64Desc() *Int64SetDesc {
	h := newInt64NodeDesc(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Int64SetDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewInt32 returns an empty skip set in ascending order.
func NewInt32() *Int32Set {
	h := newInt32Node(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Int32Set{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewInt32Desc returns an empty skip set in descending order.
func NewInt32Desc() *Int32SetDesc {
	h := newInt32NodeDesc(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Int32SetDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUint64 returns an empty skip set in ascending order.
func NewUint64() *Uint64Set {
	h := newUint64Node(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Uint64Set{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUint64Desc returns an empty skip set in descending order.
func NewUint64Desc() *Uint64SetDesc {
	h := newUint64NodeDesc(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Uint64SetDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUint32 returns an empty skip set in ascending order.
func NewUint32() *Uint32Set {
	h := newUint32Node(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Uint32Set{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUint32Desc returns an empty skip set in descending order.
func NewUint32Desc() *Uint32SetDesc {
	h := newUint32NodeDesc(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Uint32SetDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUint returns an empty skip set in ascending order.
func NewUint() *UintSet {
	h := newUintNode(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &UintSet{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewUintDesc returns an empty skip set in descending order.
func NewUintDesc() *UintSetDesc {
	h := newUintNodeDesc(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &UintSetDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}
