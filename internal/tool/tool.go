package tool

import "github.com/bilibili/gengine/internal/base"

func BinarySearch(re []*base.RuleEntity, salience int64) (int, int) {
	low := 0
	high := len(re) - 1
	mid := 0
	for low <= high {
		mid := (low + high) / 2
		if re[mid].Salience == salience {
			return low, mid
		}
		if re[mid].Salience < salience {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return low, mid
}
