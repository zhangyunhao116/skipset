package skipset

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

// NewFloat32 returns an empty skip set in ascending order.
func NewFloat32() *Float32Set {
	h := newFloat32Node(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Float32Set{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewFloat32Desc returns an empty skip set in descending order.
func NewFloat32Desc() *Float32SetDesc {
	h := newFloat32NodeDesc(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Float32SetDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewFloat64 returns an empty skip set in ascending order.
func NewFloat64() *Float64Set {
	h := newFloat64Node(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Float64Set{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
}

// NewFloat64Desc returns an empty skip set in descending order.
func NewFloat64Desc() *Float64SetDesc {
	h := newFloat64NodeDesc(0, maxLevel)
	h.flags.SetTrue(fullyLinked)
	return &Float64SetDesc{
		header:       h,
		highestLevel: defaultHighestLevel,
	}
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
