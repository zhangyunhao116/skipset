![LOGO](https://raw.githubusercontent.com/ZYunH/public-data/master/skipset-logo.png)

## Introduction

skipset is a high-performance concurrent set based on skip list. In typical pattern(one million operations, 90%CONTAINS 9%INSERT 1%DELETE), the skipset up to 3x ~ 10x faster than the built-in sync.Map.

The main idea behind the skipset is [A Simple Optimistic Skiplist Algorithm](<https://people.csail.mit.edu/shanir/publications/LazySkipList.pdf>).

Different from the sync.Map, the items in the skipset are always sorted, and the `Contains` and `Range` operations are wait-free (A goroutine is guaranteed to complete a operation as long as it keeps taking steps, regardless of the activity of other goroutines).



## Feature

- Concurrent safe API with high-performance.
- Wait-free Contains and Range methods.
- Sorted items.



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

l.Range(func(i int, score int) bool {
	fmt.Println("skipset range found ", score)
	return true
})

l.Delete(15)
fmt.Printf("skipset contains %d items\r\n", l.Len())
```



## Benchmark

Go version: go1.15 linux/amd64

CPU: AMD 3700x(8C16T), running at 3.6GHz

OS: ubuntu 18.04

MEMORY: 32GB

```shell
$ go test -run=NOTEST -bench=. -count=20 -timeout=60m -benchmem > x.txt
$ benchstat x.txt
```

```
name                                                time/op
Insert/skipset-16                                    326ns ± 8%
Insert/sync.Map-16                                   662ns ±11%
Contains100Hits/skipset-16                          9.37ns ± 4%
Contains100Hits/sync.Map-16                         6.31ns ±16%
Contains50Hits/skipset-16                           9.76ns ± 6%
Contains50Hits/sync.Map-16                          5.51ns ±11%
ContainsNoHits/skipset-16                           9.85ns ± 5%
ContainsNoHits/sync.Map-16                          4.57ns ±12%
50Insert50Contains/skipset-16                        226ns ± 6%
50Insert50Contains/sync.Map-16                       609ns ± 7%
30Insert70Contains/skipset-16                        186ns ± 6%
30Insert70Contains/sync.Map-16                       624ns ± 6%
1Delete9Insert90Contains/skipset-16                 54.7ns ± 5%
1Delete9Insert90Contains/sync.Map-16                 493ns ± 1%
1Range9Delete90Insert900Contains/skipset-16         2.71µs ± 9%
1Range9Delete90Insert900Contains/sync.Map-16        6.62µs ± 6%
StringInsert/skipset-16                              322ns ± 7%
StringInsert/sync.Map-16                            1.05µs ± 4%
StringContains50Hits/skipset-16                     19.6ns ± 5%
StringContains50Hits/sync.Map-16                    9.53ns ± 3%
String30Insert70Contains/skipset-16                  196ns ± 4%
String30Insert70Contains/sync.Map-16                 824ns ± 6%
String1Delete9Insert90Contains/skipset-16           59.0ns ± 4%
String1Delete9Insert90Contains/sync.Map-16           580ns ± 1%
String1Range9Delete90Insert900Contains/skipset-16   2.33µs ±10%
String1Range9Delete90Insert900Contains/sync.Map-16  7.68µs ± 5%

name                                                alloc/op
Insert/skipset-16                                    59.0B ± 0%
Insert/sync.Map-16                                    154B ±25%
Contains100Hits/skipset-16                           0.00B     
Contains100Hits/sync.Map-16                          0.00B     
Contains50Hits/skipset-16                            0.00B     
Contains50Hits/sync.Map-16                           0.00B     
ContainsNoHits/skipset-16                            0.00B     
ContainsNoHits/sync.Map-16                           0.00B     
50Insert50Contains/skipset-16                        29.0B ± 0%
50Insert50Contains/sync.Map-16                       88.5B ± 9%
30Insert70Contains/skipset-16                        17.0B ± 0%
30Insert70Contains/sync.Map-16                       71.1B ±11%
1Delete9Insert90Contains/skipset-16                  0.00B     
1Delete9Insert90Contains/sync.Map-16                 48.0B ± 0%
1Range9Delete90Insert900Contains/skipset-16          5.00B ± 0%
1Range9Delete90Insert900Contains/sync.Map-16        2.21kB ± 6%
StringInsert/skipset-16                              90.0B ± 0%
StringInsert/sync.Map-16                              195B ± 0%
StringContains50Hits/skipset-16                      3.00B ± 0%
StringContains50Hits/sync.Map-16                     3.00B ± 0%
String30Insert70Contains/skipset-16                  38.0B ± 0%
String30Insert70Contains/sync.Map-16                  100B ±13%
String1Delete9Insert90Contains/skipset-16            16.0B ± 0%
String1Delete9Insert90Contains/sync.Map-16           63.6B ± 1%
String1Range9Delete90Insert900Contains/skipset-16    22.0B ± 0%
String1Range9Delete90Insert900Contains/sync.Map-16  2.25kB ± 5%
```