package skipset

// NewInt return an empty int skip set in ascending order.
func NewInt() *IntSet {
	h := newIntNode(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &IntSet{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewIntDesc return an empty int skip set in descending order.
func NewIntDesc() *IntSetDesc {
	h := newIntNodeDesc(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &IntSetDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}
