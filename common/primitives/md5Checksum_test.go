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
		err := n.UnmarshalBinary(data)
		if err != nil {
			t.Error(err)
		}

		if !h.IsSameAs(n) {
			t.Error("Failed, should be same")
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
		t.Error("Failed, should be all 0s")
	}
}
