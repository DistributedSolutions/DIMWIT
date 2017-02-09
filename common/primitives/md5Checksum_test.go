package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestMd5(t *testing.T) {
	for i := 0; i < 1000; i++ {
		h := RandomMD5()
		data, _ := h.MarshalBinary()

		n := new(MD5Checksum)
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

		if h.Empty() {
			t.Error("Should not be empty")
		}
	}

	m := new(MD5Checksum)

	i, err := HexToMD5Checksum("")
	if err == nil {
		t.Error("Should fail")
	}

	str, _ := random.RandomHexStringOfSize(m.Length() * 2)
	i, err = HexToMD5Checksum(str)
	if err != nil {
		t.Error(err)
	}

	if i.String() != str {
		t.Error("Failed")
	}
}

func TestMd5Diff(t *testing.T) {
	same := 0
	for i := 0; i < 1000; i++ {
		a := RandomMD5()
		b := RandomMD5()
		if a.IsSameAs(b) {
			same++
		}
	}
	if same > 15 {
		t.Error("More than 15 are the same, it is totally random, so it is likely the IsSameAs is broken.")
	}
}

func TestEmptyMD(t *testing.T) {
	ih := new(MD5Checksum)
	if !ih.Empty() {
		t.Error("Should be empty")
	}
}

func TestBadUnmarshalMD5(t *testing.T) {
	badData := []byte{}

	n := new(Hash)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
