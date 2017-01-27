package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
)

var _ = fmt.Sprintf("")

func TestInfoHash(t *testing.T) {
	for i := 0; i < 1000; i++ {
		h := RandomInfoHash()
		data, _ := h.MarshalBinary()

		n := new(InfoHash)
		err := n.UnmarshalBinary(data)
		if err != nil {
			t.Error(err)
		}

		if !h.IsSameAs(n) {
			t.Error("Failed, should be same")
		}
	}

	i, err := HexToInfoHash("")
	if err == nil {
		t.Error("Should fail")
	}

	i, err = HexToInfoHash("0000000000000000000000000000000000000000")
	if err != nil {
		t.Error(err)
	}

	if i.String() != "0000000000000000000000000000000000000000" {
		t.Error("Failed, should be all 0s")
	}
}
