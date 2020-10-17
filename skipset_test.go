package skipset

import (
	"fmt"
	"math"
	"sync"
	"testing"
)

const randN = 1 << 20

func Example() {
	l := NewInt()

	for _, v := range []int{10, 12, 15} {
		if l.Insert(v) {
			fmt.Println("skipset insert", v)
		}
	}

	if l.Contains(10) {
		fmt.Println("skipset contains 10")
	}

	l.Range(func(i int, score int) bool {
		fmt.Println("skipset range found ", score)
		return true
	})

	l.Delete(15)
	fmt.Printf("skipset contains %d items\r\n", l.Len())
}

func TestNewInt64(t *testing.T) {
	// Correctness.
	l := NewInt64()
	if l.length != 0 {
		t.Fatal("invalid length")
	}
	if l.Contains(0) {
		t.Fatal("invalid contains")
	}

	if !l.Insert(0) || l.length != 1 {
		t.Fatal("invalid insert")
	}
	if !l.Contains(0) {
		t.Fatal("invalid contains")
	}
	if !l.Delete(0) || l.length != 0 {
		t.Fatal("invalid delete")
	}

	if !l.Insert(20) || l.length != 1 {
		t.Fatal("invalid insert")
	}
	if !l.Insert(22) || l.length != 2 {
		t.Fatal("invalid insert")
	}
	if !l.Insert(21) || l.length != 3 {
		t.Fatal("invalid insert")
	}

	l.Range(func(i int, score int64) bool {
		if i == 0 && score != 20 {
			t.Fatal("invalid range")
		}
		if i == 1 && score != 21 {
			t.Fatal("invalid range")
		}
		if i == 2 && score != 22 {
			t.Fatal("invalid range")
		}
		return true
	})

	if !l.Delete(21) || l.length != 2 {
		t.Fatal("invalid delete")
	}

	l.Range(func(i int, score int64) bool {
		if i == 0 && score != 20 {
			t.Fatal("invalid range")
		}
		if i == 1 && score != 22 {
			t.Fatal("invalid range")
		}
		return true
	})

	const num = math.MaxInt16
	// Make rand shuffle array.
	// The tastArray contains [1,num]
	testArray := make([]int64, num)
	testArray[0] = num + 1
	for i := 1; i < num; i++ {
		// We left 0, because it is the default score for head and tail.
		// If we check the skipset contains 0, there must be something wrong.
		testArray[i] = int64(i)
	}
	for i := len(testArray) - 1; i > 0; i-- { // Fisher–Yates shuffle
		j := fastrandn(uint32(i + 1))
		testArray[i], testArray[j] = testArray[j], testArray[i]
	}

	// Concurrent insert.
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		i := i
		wg.Add(1)
		go func() {
			l.Insert(testArray[i])
			wg.Done()
		}()
	}
	wg.Wait()
	if l.length != int64(num) {
		t.Fatalf("invalid length expected %d, got %d", num, l.length)
	}

	// Don't contains 0 after concurrent insertion.
	if l.Contains(0) {
		t.Fatal("contains 0 after concurrent insertion")
	}

	// Concurrent contains.
	for i := 0; i < num; i++ {
		i := i
		wg.Add(1)
		go func() {
			if !l.Contains(testArray[i]) {
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
			if !l.Delete(testArray[i]) {
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
	const smallRndN = 1 << 8
	for i := 0; i < 1<<16; i++ {
		wg.Add(1)
		go func() {
			r := fastrandn(num)
			if r == 0 { // we don't want to insert 0.
				r = 1
			}
			if r < 333 {
				l.Insert(int64(fastrandn(smallRndN)))
			} else if r < 666 {
				l.Contains(int64(fastrandn(smallRndN)))
			} else if r != 999 {
				l.Delete(int64(fastrandn(smallRndN)))
			} else {
				l.Range(func(i int, score int64) bool {
					if score == 0 { // default header and tail score
						t.Fatal("invalid content")
					}
					return true
				})
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkInsert_SkipSet(b *testing.B) {
	l := NewInt64()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Insert(int64(fastrandn(randN)))
		}
	})
}

func BenchmarkInsert_SyncMap(b *testing.B) {
	var l sync.Map
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Store(fastrandn(randN), nil)
		}
	})
}

func Benchmark50Insert50Contains_SkipSet(b *testing.B) {
	l := NewInt64()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u := fastrandn(2)
			if u == 0 {
				l.Insert(int64(fastrandn(randN)))
			} else {
				l.Contains(int64(fastrandn(randN)))
			}
		}
	})
}

func Benchmark50Insert50Contains_SyncMap(b *testing.B) {
	var l sync.Map
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u := fastrandn(2)
			if u == 0 {
				l.Store(fastrandn(randN), nil)
			} else {
				l.Load(fastrandn(randN))
			}
		}
	})
}

func Benchmark30Insert70Contains_SkipSet(b *testing.B) {
	l := NewInt64()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u := fastrandn(10)
			if u < 3 {
				l.Insert(int64(fastrandn(randN)))
			} else {
				l.Contains(int64(fastrandn(randN)))
			}
		}
	})
}

func Benchmark30Insert70Contains_SyncMap(b *testing.B) {
	var l sync.Map
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u := fastrandn(10)
			if u < 3 {
				l.Store(fastrandn(randN), nil)
			} else {
				l.Load(fastrandn(randN))
			}
		}
	})
}

func Benchmark1Delete9Insert90Contains_SkipSet(b *testing.B) {
	l := NewInt64()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u := fastrandn(100)
			if u == 1 {
				l.Insert(int64(fastrandn(randN)))
			} else if u == 2 {
				l.Delete(int64(fastrandn(randN)))
			} else {
				l.Contains(int64(fastrandn(randN)))
			}
		}
	})
}

func Benchmark1Delete9Insert90Contains_SyncMap(b *testing.B) {
	var l sync.Map
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u := fastrandn(100)
			if u == 1 {
				l.Store(fastrandn(randN), nil)
			} else if u == 2 {
				l.Delete(fastrandn(randN))
			} else {
				l.Load(fastrandn(randN))
			}
		}
	})
}

func Benchmark1Range9Delete90Insert900Contains_SkipSet(b *testing.B) {
	l := NewInt64()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u := fastrandn(1000)
			if u == 0 {
				l.Range(func(i int, score int64) bool {
					return true
				})
			} else if u > 10 && u < 20 {
				l.Delete(int64(fastrandn(randN)))
			} else if u >= 100 && u < 190 {
				l.Insert(int64(fastrandn(randN)))
			} else {
				l.Contains(int64(fastrandn(randN)))
			}
		}
	})
}

func Benchmark1Range9Delete90Insert900Contains_SyncMap(b *testing.B) {
	var l sync.Map
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			u := fastrandn(1000)
			if u == 0 {
				l.Range(func(key, value interface{}) bool {
					return true
				})
			} else if u > 10 && u < 20 {
				l.Delete(fastrandn(randN))
			} else if u >= 100 && u < 190 {
				l.Store(fastrandn(randN), nil)
			} else {
				l.Load(fastrandn(randN))
			}
		}
	})
}
