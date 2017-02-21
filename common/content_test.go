package common_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common"
	// "github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestContent(t *testing.T) {
	for i := 0; i < 1000; i++ {
		l := RandomNewContent()
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(Content)
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

	for i := 0; i < 3; i++ {
		l := RandomContentList(100)
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(ContentList)
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

func TestBadUnmarshalContent(t *testing.T) {
	badData := []byte{}

	n := new(Content)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
