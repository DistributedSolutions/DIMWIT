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
		newData, err := n.UnmarshalBinaryData(data)
		if err != nil {
			t.Error(err)
		}
		if !n.IsSameAs(l) {
			t.Error("Should match.")
		}

		if len(newData) != 0 {
			t.Error("Failed, should have no bytes left")
		}
	}
}

func TestTagList(t *testing.T) {
	for i := 0; i < 10000; i++ {
		l := RandomTagList(random.RandomUInt32Between(0, 100))
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(i, err)
		}

		n := new(TagList)
		newData, err := n.UnmarshalBinaryData(data)
		if err != nil {
			t.Error(i, err)
		}

		if !n.IsSameAs(l) {
			t.Error("Should match.")
		}
		if len(newData) != 0 {
			t.Error("Failed, should have no bytes left")
		}
	}
}

func TestBadUnmarshalTag(t *testing.T) {
	badData := []byte{}

	n := new(Tag)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}

	s := new(TagList)
	_, err = s.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
