package skipset

import (
	"testing"
)

func TestFlag(t *testing.T) {
	// Correctness.
	const (
		f0 = 1 << iota
		f1
		f2
		f3
		f4
		f5
		f6
		f7
	)
	x := &bitflag{}

	x.SetTrue(f0 | f1 | f2)
	if !x.Get(f0) || !x.Get(f1) || !x.Get(f2) || !x.MGet(f0|f1|f2, f0|f1|f2) {
		t.Fatal("invalid")
	}

	x.SetFalse(f1 | f2)
	if !x.Get(f0) || x.Get(f1) || x.Get(f2) || !x.MGet(f0|f1|f2, f0) {
		t.Fatal("invalid")
	}

	x.SetTrue(f3 | f4)
	if !x.Get(f0) || x.Get(f1) || x.Get(f2) || !x.Get(f3) || !x.Get(f4) || !x.MGet(f0|f1|f2|f3|f4, f0|f3|f4) {
		t.Fatal("invalid")
	}
}
