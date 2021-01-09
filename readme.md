<p align="center">
  <img src="https://raw.githubusercontent.com/ZYunH/public-data/master/skipset-logo2.png"/>
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

- Insert a batch of elements and then use only `Contains` (use built-in map is even better)



## QuickStart

See [Go doc](https://godoc.org/github.com/ZYunH/skipset) for more information.

```go
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
```



## Benchmark

Go version: go1.15.6 linux/amd64

CPU: AMD 3700x(8C16T), running at 3.6GHz

OS: ubuntu 18.04

MEMORY: 16G x 2 (3200MHz)

![benchmark](https://raw.githubusercontent.com/ZYunH/public-data/master/skipset-benchmark.png)

```shell
$ go test -run=NOTEST -bench=. -benchtime=100000x -benchmem -count=10 -timeout=60m  > x.txt
$ benchstat x.txt
```

```
name                                                time/op
Insert/skipset-16                                    186ns ± 9%
Insert/sync.Map-16                                   660ns ± 6%
Contains100Hits/skipset-16                          14.7ns ±12%
Contains100Hits/sync.Map-16                         14.8ns ±15%
Contains50Hits/skipset-16                           14.9ns ±10%
Contains50Hits/sync.Map-16                          13.4ns ±12%
ContainsNoHits/skipset-16                           14.4ns ± 6%
ContainsNoHits/sync.Map-16                          12.2ns ± 7%
50Insert50Contains/skipset-16                       86.7ns ± 2%
50Insert50Contains/sync.Map-16                       564ns ± 5%
30Insert70Contains/skipset-16                       66.3ns ± 4%
30Insert70Contains/sync.Map-16                       601ns ± 6%
1Delete9Insert90Contains/skipset-16                 28.0ns ±29%
1Delete9Insert90Contains/sync.Map-16                 366ns ± 6%
1Range9Delete90Insert900Contains/skipset-16         45.6ns ±21%
1Range9Delete90Insert900Contains/sync.Map-16        1.20µs ±13%
StringInsert/skipset-16                              205ns ±15%
StringInsert/sync.Map-16                             867ns ± 1%
StringContains50Hits/skipset-16                     27.2ns ±15%
StringContains50Hits/sync.Map-16                    18.6ns ±12%
String30Insert70Contains/skipset-16                 90.9ns ± 4%
String30Insert70Contains/sync.Map-16                 751ns ± 4%
String1Delete9Insert90Contains/skipset-16           50.2ns ± 7%
String1Delete9Insert90Contains/sync.Map-16           613ns ± 3%
String1Range9Delete90Insert900Contains/skipset-16   56.2ns ± 5%
String1Range9Delete90Insert900Contains/sync.Map-16  1.36µs ±14%

name                                                alloc/op
Insert/skipset-16                                    59.0B ± 0%
Insert/sync.Map-16                                    128B ± 0%
Contains100Hits/skipset-16                           0.00B     
Contains100Hits/sync.Map-16                          0.00B     
Contains50Hits/skipset-16                            0.00B     
Contains50Hits/sync.Map-16                           0.00B     
ContainsNoHits/skipset-16                            0.00B     
ContainsNoHits/sync.Map-16                           0.00B     
50Insert50Contains/skipset-16                        29.0B ± 0%
50Insert50Contains/sync.Map-16                       65.4B ± 5%
30Insert70Contains/skipset-16                        17.0B ± 0%
30Insert70Contains/sync.Map-16                       79.6B ±13%
1Delete9Insert90Contains/skipset-16                  0.00B     
1Delete9Insert90Contains/sync.Map-16                 42.9B ± 4%
1Range9Delete90Insert900Contains/skipset-16          5.00B ± 0%
1Range9Delete90Insert900Contains/sync.Map-16          300B ±15%
StringInsert/skipset-16                              90.0B ± 0%
StringInsert/sync.Map-16                              152B ± 0%
StringContains50Hits/skipset-16                      3.00B ± 0%
StringContains50Hits/sync.Map-16                     3.00B ± 0%
String30Insert70Contains/skipset-16                  38.0B ± 0%
String30Insert70Contains/sync.Map-16                 98.6B ±11%
String1Delete9Insert90Contains/skipset-16            22.0B ± 0%
String1Delete9Insert90Contains/sync.Map-16           72.6B ± 2%
String1Range9Delete90Insert900Contains/skipset-16    22.0B ± 0%
String1Range9Delete90Insert900Contains/sync.Map-16    302B ±14%
```