package objects_test

import (
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/constructor/objects"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

func TestRemoveFromList(t *testing.T) {
	list := make([]*lite.EntryHolder, 0)
	list = RemoveFromList(list, 100)
	list = append(list, new(lite.EntryHolder))
	list = append(list, new(lite.EntryHolder))
	list = append(list, new(lite.EntryHolder))
	list = RemoveFromList(list, 1)
	if len(list) != 2 {
		t.Error("List bad length")
	}
}
