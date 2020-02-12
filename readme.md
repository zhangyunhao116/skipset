![LOGO](https://raw.githubusercontent.com/ZYunH/public-data/master/skipset-logo.png)

## Introduction

skipset is a high-performance concurrent set based on skip list. In typical pattern(one million operations, 90%CONTAINS 9%INSERT 1%DELETE), the skipset up to 2x ~ 2.5x faster than the built-in sync.Map.

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

Background:  4 core CPU, clocked at 3.1 GHz. Create 8 goroutines to execute these operations.

##### 1,00,000 operations

- **90%CONTAINS 9%INSERT 1%DELETE**: 2.5x faster than sync.Map.
- **30% INSERT 70%CONTAINS**: 3.2x faster than sync.Map.
- **100% INSERT**: 1.2x ~ 1.3x faster than sync.Map, reduce memory consumption by about 50%.
- **100% RANGE(with 1,000 items)**: 3x ~ 3.5x faster than sync.Map. 
- **100% CONTAINS**: 1.3x faster than sync.Map.
- **100% DELETE**: 1.2x faster than sync.Map.

##### 1,000,000 operations

- **90%CONTAINS 9%INSERT 1%DELETE**: 2x ~ 2.5x faster than sync.Map.
- **30% INSERT 70%CONTAINS**: 2x ~ 2.1x faster than sync.Map.
- **100% INSERT**: 1.1x ~ 1.3x faster than sync.Map, reduce memory consumption by about 50%.
- **100% RANGE(with 1,000 items)**: 3x ~ 3.5x faster than sync.Map. 
- **100% CONTAINS**: equal to sync.Map.
- **100% DELETE**: equal to sync.Map.