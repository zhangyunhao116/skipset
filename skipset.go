// Package skipset is a high-performance concurrent set based on skip list.
// In typical pattern(one million operations, 90%CONTAINS 9%INSERT 1%DELETE),
// the skipset up to 3x ~ 10x faster than the built-in sync.Map.
package skipset

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// Int64Set represents a set based on skip list in ascending order.
type Int64Set struct {
	header *int64Node
	tail   *int64Node
	length int64
}

type int64Node struct {
	score int64
	next  []*int64Node
	mu    sync.Mutex
	flags bitflag
}

func newInt64Node(score int64, level int) *int64Node {
	return &int64Node{
		score: score,
		next:  make([]*int64Node, level),
	}
}

// loadNext return `n.next[i]`(atomic)
func (n *int64Node) loadNext(i int) *int64Node {
	return (*int64Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&n.next[i]))))
}

// storeNext same with `n.next[i] = val`(atomic)
func (n *int64Node) storeNext(i int, val *int64Node) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&n.next[i])), unsafe.Pointer(val))
}

// NewInt64 return an empty int64 skip set.
func NewInt64() *Int64Set {
	h, t := newInt64Node(0, maxLevel), newInt64Node(0, maxLevel)
	for i := 0; i < maxLevel; i++ {
		h.next[i] = t
	}
	h.flags.SetTrue(fullyLinked)
	t.flags.SetTrue(fullyLinked)
	return &Int64Set{
		header: h,
		tail:   t,
	}
}

// findNodeDelete takes a score and two maximal-height arrays then searches exactly as in a sequential skip-list.
// The returned preds and succs always satisfy preds[i] > score >= succs[i].
func (s *Int64Set) findNodeDelete(score int64, preds *[maxLevel]*int64Node, succs *[maxLevel]*int64Node) int {
	// lFound represents the index of the first layer at which it found a node.
	lFound, x := -1, s.header
	for i := maxLevel - 1; i >= 0; i-- {
		succ := x.loadNext(i)
		for succ != s.tail && succ.score < score {
			x = succ
			succ = x.loadNext(i)
		}
		preds[i] = x
		succs[i] = succ

		// Check if the score already in the skip list.
		if lFound == -1 && succ != s.tail && score == succ.score {
			lFound = i
		}
	}
	return lFound
}

// findNodeInsert takes a score and two maximal-height arrays then searches exactly as in a sequential skip-set.
// The returned preds and succs always satisfy preds[i] > score > succs[i].
func (s *Int64Set) findNodeInsert(score int64, preds *[maxLevel]*int64Node, succs *[maxLevel]*int64Node) int {
	x := s.header
	for i := maxLevel - 1; i >= 0; i-- {
		succ := x.loadNext(i)
		for succ != s.tail && succ.score < score {
			x = succ
			succ = x.loadNext(i)
		}
		preds[i] = x
		succs[i] = succ

		// Check if the score already in the skip list.
		if succ != s.tail && score == succ.score {
			return i
		}
	}
	return -1
}

func unlockInt64(preds [maxLevel]*int64Node, highestLevel int) {
	var prevPred *int64Node
	for i := highestLevel; i >= 0; i-- {
		if preds[i] != prevPred { // the node could be unlocked by previous loop
			preds[i].mu.Unlock()
			prevPred = preds[i]
		}
	}
}

// Insert insert the score into skip set, return true if this process insert the score into skip set,
// return false if this process can't insert this score, because another process has insert the same score.
//
// If the score is in the skip set but not fully linked, this process will wait until it is.
func (s *Int64Set) Insert(score int64) bool {
	level := randomLevel()
	var preds, succs [maxLevel]*int64Node
	for {
		lFound := s.findNodeInsert(score, &preds, &succs)
		if lFound != -1 { // indicating the score is already in the skip-list
			nodeFound := succs[lFound]
			if !nodeFound.flags.Get(marked) {
				for !nodeFound.flags.Get(fullyLinked) {
					// The node is not yet fully linked, just waits until it is.
				}
				return false
			}
			// If the node is marked, represents some other thread is in the process of deleting this node,
			// we need to add this node in next loop.
			continue
		}

		// Add this node into skip list.
		var (
			highestLocked        = -1 // the highest level being locked by this process
			valid                = true
			pred, succ, prevPred *int64Node
		)
		for layer := 0; valid && layer < level; layer++ {
			pred = preds[layer]   // target node's previous node
			succ = succs[layer]   // target node's next node
			if pred != prevPred { // the node in this layer could be locked by previous loop
				pred.mu.Lock()
				highestLocked = layer
				prevPred = pred
			}
			// valid check if there is another node has inserted into the skip list in this layer during this process.
			// It is valid if:
			// 1. The previous node and next node both are not marked.
			// 2. The previous node's next node is succ in this layer.
			valid = !pred.flags.Get(marked) && !succ.flags.Get(marked) && pred.loadNext(layer) == succ
		}
		if !valid {
			unlockInt64(preds, highestLocked)
			continue
		}

		nn := newInt64Node(score, level)
		for layer := 0; layer < level; layer++ {
			nn.next[layer] = succs[layer]
			preds[layer].storeNext(layer, nn)
		}
		nn.flags.SetTrue(fullyLinked)
		unlockInt64(preds, highestLocked)
		atomic.AddInt64(&s.length, 1)
		return true
	}
}

// Contains check if the score is in the skip set.
func (s *Int64Set) Contains(score int64) bool {
	x := s.header
	for i := maxLevel - 1; i >= 0; i-- {
		nex := x.loadNext(i)
		for nex != s.tail && nex.score < score {
			x = nex
			nex = x.loadNext(i)
		}

		// Check if the score already in the skip list.
		if nex != s.tail && score == nex.score {
			return nex.flags.MGet(fullyLinked|marked, fullyLinked)
		}
	}
	return false
}

// Delete a node from the skip set.
func (s *Int64Set) Delete(score int64) bool {
	var (
		nodeToDelete *int64Node
		isMarked     bool // represents if this operation mark the node
		topLayer     = -1
		preds, succs [maxLevel]*int64Node
	)
	for {
		lFound := s.findNodeDelete(score, &preds, &succs)
		if isMarked || // this process mark this node or we can find this node in the skip list
			lFound != -1 && succs[lFound].flags.MGet(fullyLinked|marked, fullyLinked) && (len(succs[lFound].next)-1) == lFound {
			if !isMarked { // we don't mark this node for now
				nodeToDelete = succs[lFound]
				topLayer = lFound
				nodeToDelete.mu.Lock()
				if nodeToDelete.flags.Get(marked) {
					// The node is marked by another process,
					// the physical deletion will be accomplished by another process.
					nodeToDelete.mu.Unlock()
					return false
				}
				nodeToDelete.flags.SetTrue(marked)
				isMarked = true
			}
			// Accomplish the physical deletion.
			var (
				highestLocked        = -1 // the highest level being locked by this process
				valid                = true
				pred, succ, prevPred *int64Node
			)
			for layer := 0; valid && (layer <= topLayer); layer++ {
				pred, succ = preds[layer], succs[layer]
				if pred != prevPred { // the node in this layer could be locked by previous loop
					pred.mu.Lock()
					highestLocked = layer
					prevPred = pred
				}
				// valid check if there is another node has inserted into the skip list in this layer
				// during this process, or the previous is deleted by another process.
				// It is valid if:
				// 1. the previous node exists.
				// 2. no another node has inserted into the skip list in this layer.
				valid = !pred.flags.Get(marked) && pred.loadNext(layer) == succ
			}
			if !valid {
				unlockInt64(preds, highestLocked)
				continue
			}
			for i := topLayer; i >= 0; i-- {
				// Now we own the `nodeToDelete`, no other goroutine will modify it.
				// So we don't need `nodeToDelete.loadNext`
				preds[i].storeNext(i, nodeToDelete.next[i])
			}
			nodeToDelete.mu.Unlock()
			unlockInt64(preds, highestLocked)
			atomic.AddInt64(&s.length, -1)
			return true
		}
		return false
	}
}

// Range calls f sequentially for each score present in the skip set.
// If f returns false, range stops the iteration.
func (s *Int64Set) Range(f func(score int64) bool) {
	x := s.header.loadNext(0)
	for x != s.tail {
		if !x.flags.MGet(fullyLinked|marked, fullyLinked) {
			x = x.loadNext(0)
			continue
		}
		if !f(x.score) {
			break
		}
		x = x.loadNext(0)
	}
}

// Len return the length of this skip set.
// Keep in sync with types_gen.go:lengthFunction
// Special case for code generation, Must in the tail of skipset.go.
func (s *Int64Set) Len() int {
	return int(atomic.LoadInt64(&s.length))
}
