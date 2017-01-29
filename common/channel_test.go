package common_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common"
	// "github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestChannel(t *testing.T) {
	for i := 0; i < 10; i++ {
		l := RandomNewChannel()
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(Channel)
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

func TestBadUnmarshalChannel(t *testing.T) {
	badData := []byte{}

	n := new(Channel)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
