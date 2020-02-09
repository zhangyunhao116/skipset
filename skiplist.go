package skiplist

import (
	"fmt"
	"github.com/ZYunH/lockedsource"
	atom "go.uber.org/atomic"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	maxLevel = 32
	p        = 0.25
)

var rnd = rand.New(lockedsource.New(0))

type SkipList struct {
	header *Node
	tail   *Node
	length int64
}

type Node struct {
	score       int64
	next        []*Node
	marked      atom.Bool
	fullyLinked atom.Bool
	mu          sync.Mutex
}

func newNode(score int64, level int) *Node {
	return &Node{
		score:       score,
		next:        make([]*Node, level),
		marked:      *atom.NewBool(false),
		fullyLinked: *atom.NewBool(true),
	}
}

func New() *SkipList {
	h, t := newNode(0, maxLevel), newNode(0, maxLevel)
	for i := 0; i < maxLevel; i++ {
		h.next[i] = t
	}
	return &SkipList{
		header: h,
		tail:   t,
	}
}

func randomLevel() int {
	level := 1
	for rnd.Float64() < p {
		level++
	}
	if level > maxLevel {
		return maxLevel
	}
	return level
}

// findNode takes a score and two maximal-height arrays then searches exactly as in a sequential skip-list.
// The returned preds and succs always satisfy preds[i] > score > succs[i].
func (s *SkipList) findNode(score int64, preds *[maxLevel]*Node, succs *[maxLevel]*Node) int {
	// lFound represents the index of the first layer at which it found a node.
	lFound, x := -1, s.header
	for i := maxLevel - 1; i >= 0; i-- {
		succ := x.next[i]
		for succ != s.tail && succ.score < score {
			x = succ
			succ = x.next[i]
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

// findNode takes a score and two maximal-height arrays then searches exactly as in a sequential skip-list.
// The returned preds and succs always satisfy preds[i] > score > succs[i].
func (s *SkipList) findNodeSimple(score int64, preds *[maxLevel]*Node, succs *[maxLevel]*Node) int {
	// lFound represents the index of the first layer at which it found a node.
	x := s.header
	for i := maxLevel - 1; i >= 0; i-- {
		succ := x.next[i]
		for succ != s.tail && succ.score < score {
			x = succ
			succ = x.next[i]
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

func unlock(preds [maxLevel]*Node, highestLevel int) {
	var prevPred *Node
	for i := highestLevel; i >= 0; i-- {
		if preds[i] != prevPred { // the node could be unlocked by previous loop
			preds[i].mu.Unlock()
			prevPred = preds[i]
		}
	}
}

func (s *SkipList) Insert(score int64) bool {
	level := randomLevel()
	var preds, succs [maxLevel]*Node
	for {
		lFound := s.findNodeSimple(score, &preds, &succs)
		if lFound != -1 { // indicating the score is already in the skip-list
			nodeFound := succs[lFound]
			if !nodeFound.marked.Load() {
				for !nodeFound.fullyLinked.Load() {
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
			pred, succ, prevPred *Node
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
			valid = !pred.marked.Load() && !succ.marked.Load() && pred.next[layer] == succ
		}
		if !valid {
			unlock(preds, highestLocked)
			continue
		}

		nn := newNode(score, level)
		for i := 0; i < level; i++ {
			nn.next[i] = succs[i]
			preds[i].next[i] = nn
		}
		nn.fullyLinked.Store(true)
		unlock(preds, highestLocked)
		atomic.AddInt64(&s.length, 1)
		return true
	}
}

func (s *SkipList) Contains(score int64) bool {
	x := s.header
	for i := maxLevel - 1; i >= 0; i-- {
		for x.next[i] != s.tail && x.next[i].score < score {
			x = x.next[i]
		}

		// Check if the score already in the skip list.
		if x.next[i] != s.tail && score == x.next[i].score {
			return x.next[i].fullyLinked.Load() && !x.next[i].marked.Load()
		}
	}
	return false
}

func (s *SkipList) Remove(score int64) bool {
	var (
		nodeToDelete *Node
		isMarked     bool // represents if this operation mark the node
		topLayer     = -1
		preds, succs [maxLevel]*Node
	)
	for {
		lFound := s.findNode(score, &preds, &succs)
		if isMarked || // this process mark this node or we can find this node in the skip list
			lFound != -1 && succs[lFound].fullyLinked.Load() && !succs[lFound].marked.Load() && (len(succs[lFound].next)-1) == lFound {
			if !isMarked { // we don't mark this node for now
				nodeToDelete = succs[lFound]
				topLayer = lFound
				nodeToDelete.mu.Lock()
				if nodeToDelete.marked.Load() {
					// The node is marked by another process,
					// the physical deletion will be accomplished by another process.
					nodeToDelete.mu.Unlock()
					return false
				}
				nodeToDelete.marked.Store(true)
				isMarked = true
			}
			// Accomplish the physical deletion.
			var (
				highestLocked        = -1 // the highest level being locked by this process
				valid                = true
				pred, succ, prevPred *Node
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
				valid = !pred.marked.Load() && pred.next[layer] == succ
			}
			if !valid {
				unlock(preds, highestLocked)
				continue
			}
			for i := topLayer; i >= 0; i-- {
				preds[i].next[i] = nodeToDelete.next[i]
			}
			nodeToDelete.mu.Unlock()
			unlock(preds, highestLocked)
			atomic.AddInt64(&s.length, -1)
			return true
		}
		return false
	}
}

func (s *SkipList) print() {
	for i := maxLevel - 1; i >= 0; i-- {
		print(i, " ")
		x := s.header.next[i]

		for x != s.tail {
			print("[score:", x.score,
				"] -> ")
			x = x.next[i]
		}
		print("tail")
		print("\r\n")
	}
	print("\r\n")
}

func (s *SkipList) sprint() string {
	data := make([]string, 10000)
	addS := func(a ...interface{}) {
		x := fmt.Sprint(a...)
		data = append(data, x)
	}
	for i := maxLevel - 1; i >= 0; i-- {
		addS(i, " ")
		x := s.header.next[i]

		for x != s.tail {
			addS("[score:", x.score,
				"] -> ")
			x = x.next[i]
		}
		addS("tail")
		addS("\r\n")
	}
	addS("\r\n")
	return strings.Join(data, "")
}
