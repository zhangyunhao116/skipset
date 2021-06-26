<p align="center">
  <img src="https://raw.githubusercontent.com/zhangyunhao116/public-data/master/skipset-logo2.png"/>
</p>

## Introduction

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
	l := NewInt()

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



## Benchmark

Go version: go1.16.2 linux/amd64

CPU: AMD 3700x(8C16T), running at 3.6GHz

OS: ubuntu 18.04

MEMORY: 16G x 2 (3200MHz)

![benchmark](https://raw.githubusercontent.com/zhangyunhao116/public-data/master/skipset-benchmark.png)

```shell
$ go test -run=NOTEST -bench=. -benchtime=100000x -benchmem -count=20 -timeout=60m  > x.txt
$ benchstat x.txt
```

```
name                                              time/op
Int64/Add/skipset-16                               107ns ±12%
Int64/Add/sync.Map-16                              679ns ± 5%
Int64/Contains50Hits/skipset-16                   11.8ns ±18%
Int64/Contains50Hits/sync.Map-16                  14.3ns ±19%
Int64/30Add70Contains/skipset-16                  50.9ns ±18%
Int64/30Add70Contains/sync.Map-16                  604ns ± 5%
Int64/1Remove9Add90Contains/skipset-16            32.4ns ±15%
Int64/1Remove9Add90Contains/sync.Map-16            487ns ± 5%
Int64/1Range9Remove90Add900Contains/skipset-16    35.2ns ±10%
Int64/1Range9Remove90Add900Contains/sync.Map-16   1.01µs ±18%
String/Add/skipset-16                              147ns ± 7%
String/Add/sync.Map-16                             876ns ± 5%
String/Contains50Hits/skipset-16                  19.9ns ±15%
String/Contains50Hits/sync.Map-16                 19.1ns ±16%
String/30Add70Contains/skipset-16                 66.9ns ± 5%
String/30Add70Contains/sync.Map-16                 754ns ± 4%
String/1Remove9Add90Contains/skipset-16           40.7ns ±18%
String/1Remove9Add90Contains/sync.Map-16           612ns ± 5%
String/1Range9Remove90Add900Contains/skipset-16   44.0ns ±12%
String/1Range9Remove90Add900Contains/sync.Map-16  1.24µs ±14%

name                                              alloc/op
Int64/Add/skipset-16                               58.0B ± 0%
Int64/Add/sync.Map-16                               128B ± 1%
Int64/Contains50Hits/skipset-16                    0.00B     
Int64/Contains50Hits/sync.Map-16                   0.00B     
Int64/30Add70Contains/skipset-16                   17.0B ± 0%
Int64/30Add70Contains/sync.Map-16                  80.2B ±12%
Int64/1Remove9Add90Contains/skipset-16             5.00B ± 0%
Int64/1Remove9Add90Contains/sync.Map-16            57.6B ± 5%
Int64/1Range9Remove90Add900Contains/skipset-16     5.00B ± 0%
Int64/1Range9Remove90Add900Contains/sync.Map-16     258B ±26%
String/Add/skipset-16                              90.0B ± 0%
String/Add/sync.Map-16                              152B ± 0%
String/Contains50Hits/skipset-16                   15.0B ± 0%
String/Contains50Hits/sync.Map-16                  15.0B ± 0%
String/30Add70Contains/skipset-16                  38.0B ± 0%
String/30Add70Contains/sync.Map-16                 97.6B ±12%
String/1Remove9Add90Contains/skipset-16            22.0B ± 0%
String/1Remove9Add90Contains/sync.Map-16           74.0B ± 4%
String/1Range9Remove90Add900Contains/skipset-16    22.0B ± 0%
String/1Range9Remove90Add900Contains/sync.Map-16    273B ±21%

name                                              allocs/op
Int64/Add/skipset-16                                2.00 ± 0%
Int64/Add/sync.Map-16                               4.00 ± 0%
Int64/Contains50Hits/skipset-16                     0.00     
Int64/Contains50Hits/sync.Map-16                    0.00     
Int64/30Add70Contains/skipset-16                    0.00     
Int64/30Add70Contains/sync.Map-16                   1.00 ± 0%
Int64/1Remove9Add90Contains/skipset-16              0.00     
Int64/1Remove9Add90Contains/sync.Map-16             0.00     
Int64/1Range9Remove90Add900Contains/skipset-16      0.00     
Int64/1Range9Remove90Add900Contains/sync.Map-16     0.00     
String/Add/skipset-16                               3.00 ± 0%
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