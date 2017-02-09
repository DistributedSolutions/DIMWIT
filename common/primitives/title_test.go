package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
)

var _ = fmt.Sprintf("")

func TestTitle(t *testing.T) {
	for i := 0; i < 1000; i++ {
		l := RandomTitle()
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(Title)
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

func TestEmptyTitle(t *testing.T) {
	s := new(Title)
	if !s.Empty() {
		t.Error("Should be empty")
	}

}

func TestBadUnmarshalTitle(t *testing.T) {
	badData := []byte{}

	n := new(Title)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
