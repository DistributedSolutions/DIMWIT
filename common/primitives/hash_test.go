package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestHash(t *testing.T) {
	for i := 0; i < 1000; i++ {
		h := RandomHash()
		data, _ := h.MarshalBinary()

		n := new(Hash)
		newData, err := n.UnmarshalBinaryData(data)
		if err != nil {
			t.Error(err)
		}

		if !h.IsSameAs(n) {
			t.Error("Failed, should be same")
		}

		if len(newData) != 0 {
			t.Error("Failed, should have no bytes left")
		}
	}

	m := new(Hash)

	i, err := HexToHash("")
	if err == nil {
		t.Error("Should fail")
	}

	str, _ := random.RandomHexStringOfSize(m.Length() * 2)
	i, err = HexToHash(str)
	if err != nil {
		t.Error(err)
	}

	if i.String() != str {
		t.Error("Failed")
	}
}

func TestHashList(t *testing.T) {
	for i := 0; i < 100; i++ {
		max := random.RandomUInt32Between(0, 100)

		l := RandomHashList(max)

		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(i, err)
		}

		n := new(HashList)
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

func TestBadUnmarshalHash(t *testing.T) {
	badData := []byte{}

	n := new(Hash)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
