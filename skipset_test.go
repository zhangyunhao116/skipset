package skipset

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

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

	l.Range(func(score int) bool {
		fmt.Println("skipset range found ", score)
		return true
	})

	l.Delete(15)
	fmt.Printf("skipset contains %d items\r\n", l.Len())
}

func TestIntSet(t *testing.T) {
	// Correctness.
	l := NewInt()
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

	var i int
	l.Range(func(score int) bool {
		if i == 0 && score != 20 {
			t.Fatal("invalid range")
		}
		if i == 1 && score != 21 {
			t.Fatal("invalid range")
		}
		if i == 2 && score != 22 {
			t.Fatal("invalid range")
		}
		i++
		return true
	})

	if !l.Delete(21) || l.length != 2 {
		t.Fatal("invalid delete")
	}

	i = 0
	l.Range(func(score int) bool {
		if i == 0 && score != 20 {
			t.Fatal("invalid range")
		}
		if i == 1 && score != 22 {
			t.Fatal("invalid range")
		}
		i++
		return true
	})

	const num = math.MaxInt16
	// Make rand shuffle array.
	// The testArray contains [1,num]
	testArray := make([]int, num)
	testArray[0] = num + 1
	for i := 1; i < num; i++ {
		// We left 0, because it is the default score for head and tail.
		// If we check the skipset contains 0, there must be something wrong.
		testArray[i] = int(i)
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
				panic(fmt.Sprintf("insert doesn't contains %d", i))
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
				panic(fmt.Sprintf("can't delete %d", i))
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
			if r < 333 {
				l.Insert(int(fastrandn(smallRndN)) + 1)
			} else if r < 666 {
				l.Contains(int(fastrandn(smallRndN)) + 1)
			} else if r != 999 {
				l.Delete(int(fastrandn(smallRndN)) + 1)
			} else {
				var pre int
				l.Range(func(score int) bool {
					if score <= pre { // 0 is the default value for header and tail score
						panic("invalid content")
					}
					pre = score
					return true
				})
			}
			wg.Done()
		}()
	}
	wg.Wait()

	// Correctness 2.
	var (
		x     = NewInt()
		y     = NewInt()
		count = 10000
	)

	for i := 0; i < count; i++ {
		x.Insert(i)
	}

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			x.Range(func(score int) bool {
				if x.Delete(score) {
					if !y.Insert(score) {
						panic("invalid insert")
					}
				}
				return true
			})
			wg.Done()
		}()
	}
	wg.Wait()
	if x.Len() != 0 || y.Len() != count {
		t.Fatal("invalid length")
	}

	// Concurrent Insert and Delete in small zone.
	x = NewInt()
	var (
		insertcount uint64 = 0
		deletecount uint64 = 0
	)

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {
				if fastrandn(2) == 0 {
					if x.Delete(int(fastrandn(10))) {
						atomic.AddUint64(&deletecount, 1)
					}
				} else {
					if x.Insert(int(fastrandn(10))) {
						atomic.AddUint64(&insertcount, 1)
					}
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	if insertcount < deletecount {
		panic("invalid count")
	}
	if insertcount-deletecount != uint64(x.Len()) {
		panic("invalid count")
	}

	pre := -1
	x.Range(func(score int) bool {
		if score <= pre {
			panic("invalid content")
		}
		pre = score
		return true
	})
}

func TestStringSet(t *testing.T) {
	x := NewString()
	if !x.Insert("111") || x.Len() != 1 {
		t.Fatal("invalid")
	}
	if !x.Insert("222") || x.Len() != 2 {
		t.Fatal("invalid")
	}
	if x.Insert("111") || x.Len() != 2 {
		t.Fatal("invalid")
	}
	if !x.Contains("111") || !x.Contains("222") {
		t.Fatal("invalid")
	}
	if !x.Delete("111") || x.Len() != 1 {
		t.Fatal("invalid")
	}
	if !x.Delete("222") || x.Len() != 0 {
		t.Fatal("invalid")
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		i := i
		go func() {
			if !x.Insert(strconv.Itoa(i)) {
				panic("invalid")
			}
			wg.Done()
		}()
	}
	wg.Wait()

	tmp := make([]int, 0, 100)
	x.Range(func(val string) bool {
		res, _ := strconv.Atoi(val)
		tmp = append(tmp, res)
		return true
	})
	sort.Ints(tmp)
	for i := 0; i < 100; i++ {
		if i != tmp[i] {
			t.Fatal("invalid")
		}
	}
}
