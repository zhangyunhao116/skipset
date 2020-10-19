package skipset

import (
	"math"
	"sync"
	"testing"
)

func TestFlag(t *testing.T) {
	// Correctness.
	const (
		f0 = iota
		f1
		f2
		f3
		f4
		f5
		f6
		f7
	)
	x := &bitflag{}

	// false -> true
	x.Set(f2, true)
	if !x.Get(f2) || x.data != 1<<f2 {
		t.Fatalf("invalid")
	}

	// true -> true
	x.Set(f2, true)
	if !x.Get(f2) || x.data != 1<<f2 {
		t.Fatalf("invalid")
	}

	// true -> false
	x.Set(f2, false)
	if x.Get(f2) || x.data != 0 {
		t.Fatalf("invalid")
	}

	// false -> false
	x.Set(f2, false)
	if x.Get(f2) || x.data != 0 {
		t.Fatalf("invalid")
	}

	// Concurrent set.
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		i := i
		go func() {
			x.Set(i, true)
			wg.Done()
		}()
	}
	wg.Wait()
	if x.data != math.MaxUint32 {
		t.Fatal("invalid")
	}

	for i := 0; i < 32; i++ {
		if !x.Get(i) {
			t.Fatal("invalid")
		}
	}

	for i := 0; i < 32; i++ {
		wg.Add(1)
		i := i
		go func() {
			x.Set(i, false)
			wg.Done()
		}()
	}
	wg.Wait()

	if x.data != 0 {
		t.Fatal("invalid")
	}
}
