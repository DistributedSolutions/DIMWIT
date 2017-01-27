package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestSingleTags(t *testing.T) {
	for i := 0; i < 1000; i++ {
		l := RandomTag()
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(Tag)
		err = n.UnmarshalBinary(data)
		if err != nil {
			t.Error(err)
		}
		if !n.IsSameAs(l) {
			t.Error("Should match.")
		}
	}
}

func TestTagList(t *testing.T) {
	for i := 0; i < 100; i++ {
		l := RandomTagList(random.RandomUInt32Between(0, 100))
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(i, err)
		}

		n := new(TagList)
		err = n.UnmarshalBinary(data)
		if err != nil {
			t.Error(i, err)
		}

		if !n.IsSameAs(l) {
			t.Error("Should match.")
		}
	}
}
