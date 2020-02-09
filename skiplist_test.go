package skiplist

import (
	"math/rand"
	"sync"
	"testing"
)

type benchArrayCache struct {
	length      int
	itemMap     map[int64]struct{}
	Insert      []int64
	Check       []int64
	InvalidItem []int64
	Rnd37       []bool // 30% true, 70% false
	Rnd55       []bool // 50% true, 50% false
	Rnd9091     []int  // 90% ZERO, 9% ONE, 1% TWO

	count int64
	mu    sync.Mutex
}

func newBench(length int) *benchArrayCache {
	c := &benchArrayCache{}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.length = length
	c.Insert = make([]int64, length)
	c.Check = make([]int64, length)
	c.InvalidItem = make([]int64, length)
	c.Rnd37 = make([]bool, length)
	c.Rnd55 = make([]bool, length)
	c.Rnd9091 = make([]int, length)

	c.itemMap = make(map[int64]struct{}, length)
	for i := 0; i < c.length; i++ {
		c.itemMap[int64(i)] = struct{}{}
	}

	var i int
	for k := range c.itemMap {
		c.Insert[i] = k
		i++
	}
	i = 0

	for k := range c.itemMap {
		c.Check[i] = k
		i++
	}
	i = 0

	for i < length {
		c.InvalidItem[i] = int64(i + length)
		i++
	}

	for i := 0; i < length; i++ {
		nowRnd := rand.Intn(100)
		if nowRnd < 30 {
			c.Rnd37[i] = true
		}

		if nowRnd < 50 {
			c.Rnd55[i] = true
		}

		if nowRnd < 90 {
			c.Rnd9091[i] = 0
		} else if nowRnd == 99 {
			c.Rnd9091[i] = 2
		} else {
			c.Rnd9091[i] = 1
		}
	}
	return c
}

func (c *benchArrayCache) next() (n int64) {
	c.mu.Lock()
	n = c.count
	c.count++
	c.mu.Unlock()
	return n
}

func (c *benchArrayCache) rcount() {
	c.count = 0
}

var benchArray = newBench(1.1 * 10000000)

func TestNew(t *testing.T) {
	// Correctness.
	l := New()
	if l.length != 0 {
		t.Fatal("invalid length")
	}

	if !l.Insert(0) || l.length != 1 {
		t.Fatal("invalid insert")
	}
	if !l.Contains(0) {
		t.Fatal("invalid contains")
	}
	if !l.Remove(0) || l.length != 0 {
		t.Fatal("invalid delete")
	}

	// Concurrent insert.
	num := 1000000 // one million
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		i := i
		wg.Add(1)
		go func() {
			l.Insert(benchArray.Insert[i])
			wg.Done()
		}()
	}
	wg.Wait()
	if l.length != int64(num) {
		t.Fatalf("invalid length expected %d, got %d", num, l.length)
	}

	// Concurrent contains.
	for i := 0; i < num; i++ {
		i := i
		wg.Add(1)
		go func() {
			if !l.Contains(benchArray.Insert[i]) {
				wg.Done()
				t.Fatalf("insert dosen't contains %d", i)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	// Concurrent delete.
	for i := 0; i < num; i++ {
		i := i
		wg.Add(1)
		go func() {
			if !l.Remove(benchArray.Insert[i]) {
				wg.Done()
				t.Fatalf("can't delete %d", i)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if l.length != 0 {
		t.Fatalf("invalid length expected %d, got %d", 0, l.length)
	}

	// Test all methods.
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			r := rand.Intn(3)
			if r == 0 {
				l.Insert(benchArray.Insert[rand.Intn(num)])
			} else if r == 1 {
				l.Contains(benchArray.Insert[rand.Intn(num)])
			} else {
				l.Remove(benchArray.Insert[rand.Intn(num)])
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkInsert_SkipSet(b *testing.B) {
	l := New()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Insert(benchArray.Insert[benchArray.next()])
		}
	})
	benchArray.rcount()
}

func BenchmarkInsert_SyncMap(b *testing.B) {
	var l sync.Map
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Store(benchArray.Insert[benchArray.next()], nil)
		}
	})
	benchArray.rcount()
}

func Benchmark50Insert50Contains_SkipSet(b *testing.B) {
	l := New()
	var i int
	for i < 1000 {
		l.Insert(benchArray.Insert[i])
		i++
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u := benchArray.next()
			if benchArray.Rnd55[u] == true {
				l.Insert(benchArray.Insert[u])
			} else {
				l.Contains(benchArray.Check[u])
			}
		}
	})
	benchArray.rcount()
}

func Benchmark50Insert50Contains_SyncMap(b *testing.B) {
	var l sync.Map
	var i int
	for i < 1000 {
		l.Store(benchArray.Insert[i], nil)
		i++
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u := benchArray.next()
			if benchArray.Rnd55[u] == true {
				l.Store(benchArray.Insert[u], nil)
			} else {
				l.Load(benchArray.Insert[u])
			}
		}
	})
	benchArray.rcount()
}

func Benchmark30Insert70Contains_SkipSet(b *testing.B) {
	l := New()
	var i int
	for i < 1000 {
		l.Insert(benchArray.Insert[i])
		i++
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u := benchArray.next()
			if benchArray.Rnd37[u] == true {
				l.Insert(benchArray.Insert[u])
			} else {
				l.Contains(benchArray.Check[u])
			}
		}
	})
	benchArray.rcount()
}

func Benchmark30Insert70Contains_SyncMap(b *testing.B) {
	var l sync.Map
	var i int
	for i < 1000 {
		l.Store(benchArray.Insert[i], nil)
		i++
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u := benchArray.next()
			if benchArray.Rnd37[u] == true {
				l.Store(benchArray.Insert[u], nil)
			} else {
				l.Load(benchArray.Insert[u])
			}
		}
	})
	benchArray.rcount()
}

func Benchmark10Insert90Contains_SkipSet(b *testing.B) {
	l := New()
	var i int
	for i < 1000 {
		l.Insert(benchArray.Insert[i])
		i++
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u := benchArray.next()
			if benchArray.Rnd9091[u] != 0 {
				l.Insert(benchArray.Insert[u])
			} else {
				l.Contains(benchArray.Check[u])
			}
		}
	})
	benchArray.rcount()
}

func Benchmark10Insert90Contains_SyncMap(b *testing.B) {
	var l sync.Map
	var i int
	for i < 1000 {
		l.Store(benchArray.Insert[i], nil)
		i++
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u := benchArray.next()
			if benchArray.Rnd9091[u] != 0 {
				l.Store(benchArray.Insert[u], nil)
			} else {
				l.Load(benchArray.Insert[u])
			}
		}
	})
	benchArray.rcount()
}
