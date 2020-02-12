// Package skipset is a high-performance concurrent set based on skip list.
// In typical pattern(one million operations, 90%CONTAINS 9%INSERT 1%DELETE),
// the skipset up to 2x ~ 2.5x faster than the built-in sync.Map.
package skipset
