<p align="center">
  <img src="https://raw.githubusercontent.com/zhangyunhao116/public-data/master/skipset-logo2.png"/>
</p>

## Introduction

> From v0.12.0, the skipset requires Go version >= 1.18, if your Go version is lower, use v0.11.0 instead.

skipset is a high-performance, scalable, concurrent-safe set based on skip-list. In the typical pattern(100000 operations, 90%CONTAINS 9%ADD 1%REMOVE, 8C16T), the skipset up to 15x faster than the built-in `sync.Map`.

The main idea behind the skipset is [A Simple Optimistic Skiplist Algorithm](<https://people.csail.mit.edu/shanir/publications/LazySkipList.pdf>).

Different from the sync.Map, the items in the skipset are always sorted, and the `Contains` and `Range` operations are wait-free (A goroutine is guaranteed to complete an operation as long as it keeps taking steps, regardless of the activity of other goroutines).

The skipset is a set instead of a map, if you need a high-performance full replacement of `sync.Map`, see [skipmap](<https://github.com/zhangyunhao116/skipmap>).

## Features

- Scalable, high-performance, concurrent-safe.
- Wait-free Contains and Range operations (wait-free algorithms have stronger guarantees than lock-free).
- Sorted items.



## When should you use skipset

In most cases, `skipset` is better than `sync.Map`, especially in these situations: 

- **Concurrent calls multiple operations**. Such as use both `Range` and `Add` at the same time, in this situation, use skipset can obtain very large improvement on performance.
- **Memory intensive**. The skipset save at least 50% memory in the benchmark.

If only one goroutine access the set for the most of the time, such as insert a batch of elements and then use only `Contains` or `Range`, use built-in map is better.



## QuickStart

See [Go doc](https://godoc.org/github.com/zhangyunhao116/skipset) for more information.

```go
package main

import (
	"fmt"

	"github.com/zhangyunhao116/skipset"
)

func main() {
	l := skipset.NewInt()

	for _, v := range []int{10, 12, 15} {
		if l.Add(v) {
			fmt.Println("skipset add", v)
		}
	}

	if l.Contains(10) {
		fmt.Println("skipset contains 10")
	}

	l.Range(func(value int) bool {
		fmt.Println("skipset range found ", value)
		return true
	})

	l.Remove(15)
	fmt.Printf("skipset contains %d items\r\n", l.Len())
}

```



From `v0.12.0`, you can use generic version APIs.

**Note that generic APIs are always slower than typed APIs, but are more suitable for some scenarios such as functional programming.**

> e.g. `New[int]` is ~2x slower than `NewInt`, and `NewFunc(func(a, b int) bool { return a < b })` is ~2x slower than `New[int]`.
>
> Performance ranking: NewInt > New[Int] > NewFunc(func(a, b int) bool { return a < b })

```go
package main

import (
	"math"

	"github.com/zhangyunhao116/skipset"
)

func main() {
	x1 := skipset.New[int]()
	for _, v := range []int{2, 1, 3} {
		x1.Add(v)
	}
	x1.Range(func(value int) bool {
		println(value)
		return true
	})
	x2 := skipset.NewFunc(func(a, b float64) bool {
		return a < b || (math.IsNaN(a) && !math.IsNaN(b))
	})
	for _, v := range []float64{math.NaN(), 3, 1, math.NaN(), 2} {
		x2.Add(v)
	}
	x2.Range(func(value float64) bool {
		println(value)
		return true
	})
}

```

## Benchmark

> base on typed APIs.

Go version: go1.16.2 linux/amd64

CPU: AMD 3700x(8C16T), running at 3.6GHz

OS: ubuntu 18.04

MEMORY: 16G x 2 (3200MHz)

![benchmark](https://raw.githubusercontent.com/zhangyunhao116/public-data/master/skipset-benchmark.png)

```
name                                              time/op
Int64/Add/skipset-16                              86.6ns ±11%
Int64/Add/sync.Map-16                              674ns ± 6%
Int64/Contains50Hits/skipset-16                   9.85ns ±12%
Int64/Contains50Hits/sync.Map-16                  14.7ns ±30%
Int64/30Add70Contains/skipset-16                  38.8ns ±18%
Int64/30Add70Contains/sync.Map-16                  586ns ± 5%
Int64/1Remove9Add90Contains/skipset-16            24.9ns ±17%
Int64/1Remove9Add90Contains/sync.Map-16            493ns ± 5%
Int64/1Range9Remove90Add900Contains/skipset-16    25.9ns ±16%
Int64/1Range9Remove90Add900Contains/sync.Map-16   1.00µs ±12%
String/Add/skipset-16                              130ns ±14%
String/Add/sync.Map-16                             878ns ± 4%
String/Contains50Hits/skipset-16                  18.3ns ± 9%
String/Contains50Hits/sync.Map-16                 19.2ns ±18%
String/30Add70Contains/skipset-16                 61.0ns ±15%
String/30Add70Contains/sync.Map-16                 756ns ± 7%
String/1Remove9Add90Contains/skipset-16           31.3ns ±13%
String/1Remove9Add90Contains/sync.Map-16           614ns ± 6%
String/1Range9Remove90Add900Contains/skipset-16   36.2ns ±18%
String/1Range9Remove90Add900Contains/sync.Map-16  1.20µs ±17%

name                                              alloc/op
Int64/Add/skipset-16                               65.0B ± 0%
Int64/Add/sync.Map-16                               128B ± 1%
Int64/Contains50Hits/skipset-16                    0.00B     
Int64/Contains50Hits/sync.Map-16                   0.00B     
Int64/30Add70Contains/skipset-16                   19.0B ± 0%
Int64/30Add70Contains/sync.Map-16                  77.7B ±16%
Int64/1Remove9Add90Contains/skipset-16             5.00B ± 0%
Int64/1Remove9Add90Contains/sync.Map-16            57.5B ± 4%
Int64/1Range9Remove90Add900Contains/skipset-16     5.00B ± 0%
Int64/1Range9Remove90Add900Contains/sync.Map-16     255B ±22%
String/Add/skipset-16                              97.0B ± 0%
String/Add/sync.Map-16                              152B ± 0%
String/Contains50Hits/skipset-16                   15.0B ± 0%
String/Contains50Hits/sync.Map-16                  15.0B ± 0%
String/30Add70Contains/skipset-16                  40.0B ± 0%
String/30Add70Contains/sync.Map-16                 98.2B ±11%
String/1Remove9Add90Contains/skipset-16            23.0B ± 0%
String/1Remove9Add90Contains/sync.Map-16           73.9B ± 4%
String/1Range9Remove90Add900Contains/skipset-16    23.0B ± 0%
String/1Range9Remove90Add900Contains/sync.Map-16    261B ±18%

name                                              allocs/op
Int64/Add/skipset-16                                1.00 ± 0%
Int64/Add/sync.Map-16                               4.00 ± 0%
Int64/Contains50Hits/skipset-16                     0.00     
Int64/Contains50Hits/sync.Map-16                    0.00     
Int64/30Add70Contains/skipset-16                    0.00     
Int64/30Add70Contains/sync.Map-16                   1.00 ± 0%
Int64/1Remove9Add90Contains/skipset-16              0.00     
Int64/1Remove9Add90Contains/sync.Map-16             0.00     
Int64/1Range9Remove90Add900Contains/skipset-16      0.00     
Int64/1Range9Remove90Add900Contains/sync.Map-16     0.00     
String/Add/skipset-16                               2.00 ± 0%
String/Add/sync.Map-16                              5.00 ± 0%
String/Contains50Hits/skipset-16                    1.00 ± 0%
String/Contains50Hits/sync.Map-16                   1.00 ± 0%
String/30Add70Contains/skipset-16                   1.00 ± 0%
String/30Add70Contains/sync.Map-16                  2.00 ± 0%
String/1Remove9Add90Contains/skipset-16             1.00 ± 0%
String/1Remove9Add90Contains/sync.Map-16            1.00 ± 0%
String/1Range9Remove90Add900Contains/skipset-16     1.00 ± 0%
String/1Range9Remove90Add900Contains/sync.Map-16    1.00 ± 0%
```