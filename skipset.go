// Package skipset is a high-performance, scalable, concurrent-safe set based on skip-list.
// In the typical pattern(100000 operations, 90%CONTAINS 9%Add 1%Remove, 8C16T), the skipset
// up to 15x faster than the built-in sync.Map.
package skipset

// New return an empty skip set in ascending order.
func New[T ordered]() *OrderedSet[T] {
	var t T
	h := newOrderedNode(t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &OrderedSet[T]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewDesc return an empty skip set in descending order.
func NewDesc[T ordered]() *OrderedSetDesc[T] {
	var t T
	h := newOrderedNodeDesc(t, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &OrderedSetDesc[T]{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewFunc return an empty skip set in ascending order.
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
