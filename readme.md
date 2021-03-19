<p align="center">
  <img src="https://raw.githubusercontent.com/zhangyunhao116/public-data/master/skipset-logo2.png"/>
</p>

## Introduction

skipset is a high-performance concurrent set based on skip list. In typical pattern(one million operations, 90%CONTAINS 9%INSERT 1%DELETE), the skipset up to 3x ~ 10x faster than the built-in sync.Map.

The main idea behind the skipset is [A Simple Optimistic Skiplist Algorithm](<https://people.csail.mit.edu/shanir/publications/LazySkipList.pdf>).

Different from the sync.Map, the items in the skipset are always sorted, and the `Contains` and `Range` operations are wait-free (A goroutine is guaranteed to complete a operation as long as it keeps taking steps, regardless of the activity of other goroutines).



## Features

- Concurrent safe API with high-performance.
- Wait-free Contains and Range operations.
- Sorted items.



## When should you use skipset

In these situations, `skipset` is better

- **Sorted elements is needed**.
- **Concurrent calls multiple operations**. such as use both `Contains` and `Insert` at the same time.
- **Memory intensive**. The skipset save at least 50% memory in the benchmark.

In these situations, `sync.Map` is better

- Only one goroutine access the set for most of the time, such as insert a batch of elements and then use only `Contains` (use built-in map is even better).



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

```



## Benchmark

Go version: go1.16.2 linux/amd64

CPU: AMD 3700x(8C16T), running at 3.6GHz

OS: ubuntu 18.04

MEMORY: 16G x 2 (3200MHz)

![benchmark](https://raw.githubusercontent.com/zhangyunhao116/public-data/master/skipset-benchmark.png)

```shell
$ go test -run=NOTEST -bench=. -benchtime=100000x -benchmem -count=10 -timeout=60m  > x.txt
$ benchstat x.txt
```

```
name                                                time/op
Insert/skipset-16                                    107ns ±10%
Insert/sync.Map-16                                   684ns ± 5%
Contains100Hits/skipset-16                          13.1ns ±21%
Contains100Hits/sync.Map-16                         15.0ns ± 8%
Contains50Hits/skipset-16                           13.2ns ±18%
Contains50Hits/sync.Map-16                          14.0ns ±14%
ContainsNoHits/skipset-16                           14.0ns ±19%
ContainsNoHits/sync.Map-16                          12.7ns ±22%
50Insert50Contains/skipset-16                       62.8ns ± 6%
50Insert50Contains/sync.Map-16                       577ns ± 6%
30Insert70Contains/skipset-16                       44.2ns ±26%
30Insert70Contains/sync.Map-16                       583ns ± 8%
1Delete9Insert90Contains/skipset-16                 33.6ns ±13%
1Delete9Insert90Contains/sync.Map-16                 502ns ± 8%
1Range9Delete90Insert900Contains/skipset-16         37.3ns ±19%
1Range9Delete90Insert900Contains/sync.Map-16        1.13µs ±12%
StringInsert/skipset-16                              145ns ±17%
StringInsert/sync.Map-16                             877ns ± 3%
StringContains50Hits/skipset-16                     20.9ns ± 5%
StringContains50Hits/sync.Map-16                    19.6ns ±14%
String30Insert70Contains/skipset-16                 68.4ns ± 6%
String30Insert70Contains/sync.Map-16                 751ns ± 7%
String1Delete9Insert90Contains/skipset-16           39.4ns ±32%
String1Delete9Insert90Contains/sync.Map-16           617ns ± 3%
String1Range9Delete90Insert900Contains/skipset-16   44.8ns ±16%
String1Range9Delete90Insert900Contains/sync.Map-16  1.39µs ±15%

name                                                alloc/op
Insert/skipset-16                                    58.0B ± 0%
Insert/sync.Map-16                                    128B ± 0%
Contains100Hits/skipset-16                           0.00B     
Contains100Hits/sync.Map-16                          0.00B     
Contains50Hits/skipset-16                            0.00B     
Contains50Hits/sync.Map-16                           0.00B     
ContainsNoHits/skipset-16                            0.00B     
ContainsNoHits/sync.Map-16                           0.00B     
50Insert50Contains/skipset-16                        29.0B ± 0%
50Insert50Contains/sync.Map-16                       69.1B ±19%
30Insert70Contains/skipset-16                        17.0B ± 0%
30Insert70Contains/sync.Map-16                       76.9B ± 7%
1Delete9Insert90Contains/skipset-16                  5.00B ± 0%
1Delete9Insert90Contains/sync.Map-16                 55.3B ± 5%
1Range9Delete90Insert900Contains/skipset-16          5.00B ± 0%
1Range9Delete90Insert900Contains/sync.Map-16          283B ±14%
StringInsert/skipset-16                              90.0B ± 0%
StringInsert/sync.Map-16                              152B ± 0%
StringContains50Hits/skipset-16                      3.00B ± 0%
StringContains50Hits/sync.Map-16                     3.00B ± 0%
String30Insert70Contains/skipset-16                  38.0B ± 0%
String30Insert70Contains/sync.Map-16                 97.8B ±11%
String1Delete9Insert90Contains/skipset-16            22.0B ± 0%
String1Delete9Insert90Contains/sync.Map-16           72.5B ± 1%
String1Range9Delete90Insert900Contains/skipset-16    22.0B ± 0%
String1Range9Delete90Insert900Contains/sync.Map-16    315B ±22%
```