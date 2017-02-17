package objects

import (
	"time"
)

// InsideTimeWindow returns true if ele is withing window time of Main
func InsideTimeWindow(main time.Time, ele time.Time, window int64) bool {
	beg := main.Unix()
	end := ele.Unix()
	diff := end - beg
	if diff < 0 {
		diff = -1 * diff
	}
	if diff > window {
		return false
	}
	return true
}
