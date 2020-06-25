![LOGO](https://raw.githubusercontent.com/ZYunH/public-data/master/skipset-logo.png)

## Introduction

skipset is a high-performance concurrent set based on skip list. In typical pattern(one million operations, 90%CONTAINS 9%INSERT 1%DELETE), the skipset up to 3x ~ 6x faster than the built-in sync.Map.

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

**The benchmark is different on different machines, run `sh bench.sh` to get your own benchmark.(change the parameters according to your machine)**

**In most cases, the  skipset up to 3x ~ 6x faster than the built-in sync.Map in the typical pattern(one million operations, 90%CONTAINS 9%INSERT 1%DELETE)**

VERSION: go1.14 linux/amd64

CPU: 8 core CPU (Intel 9700k).

OS: ubuntu 18.04

MEMORY: 32GB

Create 8 goroutines to execute these operations.

##### 1,00,000 operations

- **90%CONTAINS 9%INSERT 1%DELETE**: 5.7x faster than sync.Map.
- **30% INSERT 70%CONTAINS**: 5.8x faster than sync.Map.
- **100% INSERT**: 3.7x faster than sync.Map, reduce memory consumption by about 50%.
- **100% RANGE(with 1,000 items)**: 2.5x faster than sync.Map. 
- **100% CONTAINS**: equal to sync.Map.
- **100% DELETE**: equal to sync.Map.

##### 1,000,000 operations

- **90%CONTAINS 9%INSERT 1%DELETE**: 5.5x faster than sync.Map.
- **30% INSERT 70%CONTAINS**: 3x faster than sync.Map.
- **100% INSERT**: 4.4x faster than sync.Map, reduce memory consumption by about 50%.
- **100% RANGE(with 1,000 items)**: 2.5x faster than sync.Map. 
- **100% CONTAINS**: equal to sync.Map.
- **100% DELETE**: equal to sync.Map.