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
name                                          time/op
Insert/skipset-16                              326ns ± 7%
Insert/sync.Map-16                             662ns ±13%
Contains100Hits/skipset-16                    10.6ns ± 4%
Contains100Hits/sync.Map-16                   6.14ns ±14%
Contains50Hits/skipset-16                     11.1ns ± 5%
Contains50Hits/sync.Map-16                    5.44ns ±11%
50Insert50Contains/skipset-16                  231ns ± 4%
50Insert50Contains/sync.Map-16                 609ns ± 6%
30Insert70Contains/skipset-16                  188ns ± 6%
30Insert70Contains/sync.Map-16                 627ns ± 5%
1Delete9Insert90Contains/skipset-16           55.6ns ± 4%
1Delete9Insert90Contains/sync.Map-16           492ns ± 1%
1Range9Delete90Insert900Contains/skipset-16   2.68µs ±10%
1Range9Delete90Insert900Contains/sync.Map-16  6.56µs ± 9%

name                                          alloc/op
Insert/skipset-16                              59.0B ± 0%
Insert/sync.Map-16                              157B ±23%
Contains100Hits/skipset-16                     0.00B     
Contains100Hits/sync.Map-16                    0.00B     
Contains50Hits/skipset-16                      0.00B     
Contains50Hits/sync.Map-16                     0.00B     
50Insert50Contains/skipset-16                  29.0B ± 0%
50Insert50Contains/sync.Map-16                 87.7B ± 7%
30Insert70Contains/skipset-16                  17.0B ± 0%
30Insert70Contains/sync.Map-16                 72.9B ±14%
1Delete9Insert90Contains/skipset-16            0.00B     
1Delete9Insert90Contains/sync.Map-16           48.0B ± 0%
1Range9Delete90Insert900Contains/skipset-16    5.00B ± 0%
1Range9Delete90Insert900Contains/sync.Map-16  2.19kB ±11%
```