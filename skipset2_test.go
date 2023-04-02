package skipset

import (
	"testing"

	"github.com/zhangyunhao116/fastrand"
)

func TestRangeFrom(t *testing.T) {
	s := NewInt64()
	sd := NewInt64Desc()
	for _, v := range []int64{-3, -1, 1, 2, 4, 6} {
		s.Add(v)
		sd.Add(v)
	}

	// s := []int64{-3, -1, 1, 2, 4, 6}
	checkRangeFrom(t, NewInt64(), fastrand.Int63(), []int64{})
	checkRangeFrom(t, s, -5, []int64{-3, -1, 1, 2, 4, 6})
	checkRangeFrom(t, s, -4, []int64{-3, -1, 1, 2, 4, 6})
	checkRangeFrom(t, s, -3, []int64{-3, -1, 1, 2, 4, 6})
	checkRangeFrom(t, s, -2, []int64{-1, 1, 2, 4, 6})
	checkRangeFrom(t, s, -1, []int64{-1, 1, 2, 4, 6})
	checkRangeFrom(t, s, 0, []int64{1, 2, 4, 6})
	checkRangeFrom(t, s, 1, []int64{1, 2, 4, 6})
	checkRangeFrom(t, s, 2, []int64{2, 4, 6})
	checkRangeFrom(t, s, 3, []int64{4, 6})
	checkRangeFrom(t, s, 4, []int64{4, 6})
	checkRangeFrom(t, s, 5, []int64{6})
	checkRangeFrom(t, s, 6, []int64{6})
	checkRangeFrom(t, s, 7, []int64{})
	checkRangeFrom(t, s, 100000, []int64{})

	// sr := []int64{6, 4, 2, 1, -1, -3}
	checkRangeFrom(t, NewInt64Desc(), fastrand.Int63(), []int64{})
	checkRangeFrom(t, sd, -5, []int64{})
	checkRangeFrom(t, sd, -4, []int64{})
	checkRangeFrom(t, sd, -3, []int64{-3})
	checkRangeFrom(t, sd, -2, []int64{-3})
	checkRangeFrom(t, sd, -1, []int64{-1, -3})
	checkRangeFrom(t, sd, 0, []int64{-1, -3})
	checkRangeFrom(t, sd, 1, []int64{1, -1, -3})
	checkRangeFrom(t, sd, 2, []int64{2, 1, -1, -3})
	checkRangeFrom(t, sd, 3, []int64{2, 1, -1, -3})
	checkRangeFrom(t, sd, 4, []int64{4, 2, 1, -1, -3})
	checkRangeFrom(t, sd, 5, []int64{4, 2, 1, -1, -3})
	checkRangeFrom(t, sd, 6, []int64{6, 4, 2, 1, -1, -3})
	checkRangeFrom(t, sd, 7, []int64{6, 4, 2, 1, -1, -3})
	checkRangeFrom(t, sd, 100000, []int64{6, 4, 2, 1, -1, -3})

	// Test wide-range values.
	s = NewInt64()
	case1len := 1000
	case1 := make([]int64, case1len)
	for i := int64(0); i < int64(case1len); i++ {
		v := fastrand.Int63()
		case1[i] = v
		s.Add(v)
	}
	insertionSort(case1)
	for i := 0; i < case1len; i++ {
		testStart := case1[i] + 1
		expected := case1[i+1:]
		checkRangeFrom(t, s, testStart, expected)
	}
}

func checkRangeFrom[T anyskipset[int64]](t *testing.T, s T, start int64, expected []int64) {
	var got []int64

	s.RangeFrom(start, func(value int64) bool {
		got = append(got, value)
		return true
	})

	if !slicesEqual(got, expected) {
		t.Fatalf("Expected: %+v (start from %v)\n Got: %+v\n", expected, start, got)
	}
}

// TODO: replace it with slices.Sort if possible.
func insertionSort[T ordered](v []T) {
	for cur := 1; cur < len(v); cur++ {
		for j := cur; j > 0 && v[j] < v[j-1]; j-- {
			v[j], v[j-1] = v[j-1], v[j]
		}
	}
}

// TODO: replace it with slices.Equal if possible.
func slicesEqual[E comparable](s1, s2 []E) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
