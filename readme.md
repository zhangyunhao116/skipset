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
$ go test -run=NOTEST -bench=. -count=20 -timeout=60m > x.txt
$ benchstat x.txt
```

```
name                                          time/op
Insert_SkipSet-16                              137ns ± 6%
Insert_SyncMap-16                              595ns ± 4%
50Insert50Contains_SkipSet-16                  123ns ± 2%
50Insert50Contains_SyncMap-16                  591ns ± 8%
30Insert70Contains_SkipSet-16                  114ns ± 2%
30Insert70Contains_SyncMap-16                  569ns ± 6%
1Delete9Insert90Contains_SkipSet-16           50.8ns ± 3%
1Delete9Insert90Contains_SyncMap-16            503ns ± 1%
1Range9Delete90Insert900Contains_SkipSet-16   1.67µs ± 8%
1Range9Delete90Insert900Contains_SyncMap-16   6.25µs ± 9%

name                                         alloc/op
Insert_SkipSet-16                             7.00B ± 0%
Insert_SyncMap-16                             91.5B ± 2%
50Insert50Contains_SkipSet-16                 5.75B ±13%
50Insert50Contains_SyncMap-16                 59.9B ±15%
30Insert70Contains_SkipSet-16                 4.00B ± 0%
30Insert70Contains_SyncMap-16                 63.1B ±16%
1Delete9Insert90Contains_SkipSet-16           0.00B     
1Delete9Insert90Contains_SyncMap-16           46.0B ± 0%
1Range9Delete90Insert900Contains_SkipSet-16   3.00B ± 0%
1Range9Delete90Insert900Contains_SyncMap-16  2.11kB ± 7%
```