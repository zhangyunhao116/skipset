package skipset

import "sync/atomic"

const (
	fullyLinked = 1 << iota
	marked
)

type bitflag struct {
	data uint32
}

func (f *bitflag) Set(flag int, val bool) {
	if val {
		for {
			old := atomic.LoadUint32(&f.data)
			if old&(1<<flag) == 0 {
				// Flag is 0, need set it to 1.
				n := old | (1 << flag)
				if atomic.CompareAndSwapUint32(&f.data, old, n) {
					return
				}
				continue
			}
			return
		}
	} else {
		for {
			old := atomic.LoadUint32(&f.data)
			if old&(1<<flag) != 0 {
				// Flag is 1, need set it to 0.
				n := old ^ (1 << flag)
				if atomic.CompareAndSwapUint32(&f.data, old, n) {
					return
				}
				continue
			}
			return
		}
	}
}

func (f *bitflag) Get(flag int) bool {
	return (atomic.LoadUint32(&f.data) & (1 << flag)) != 0
}
