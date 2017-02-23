package objects

import (
	"time"

	"github.com/DistributedSolutions/DIMWIT/factom-lite"
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

func RemoveFromList(list []*lite.EntryHolder, i int) (newList []*lite.EntryHolder) {
	defer func() {
		if r := recover(); r != nil {
			newList = list
			return
		}
	}()
	list[i] = list[len(list)-1]
	return list[:len(list)-1]
	//return append(list[:i], list[i+1:]...)
}

//append(list[:i], list[i+1:]...)
