package test

import (
	"fmt"
	"testing"
)

func Test_binary(t *testing.T) {

	for i := 110; i > 40; i-- {

		raw := []int64{100, 90, 80, 70, 70, 60, 50, 50}

		low, mid := binarySearch(raw, int64(i))
		ire := []int64{int64(i)}
		if mid == 0 {
			newRe := append(ire, raw[low:]...)
			raw = append(raw[:low], newRe...)
		} else {
			newRe := append(ire, raw[mid:]...)
			raw = append(raw[:mid], newRe...)
		}
		println(fmt.Sprintf("new raw:%+v", raw))
	}

}

func binarySearch(re []int64, salience int64) (int, int) {
	low := 0
	high := len(re) - 1
	mid := 0
	for low <= high {
		mid := (low + high) / 2
		if re[mid] == salience {
			return low, mid
		}
		if re[mid] < salience {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return low, mid
}
