package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
)

var _ = fmt.Sprintf("")

func TestHash(t *testing.T) {
	for i := 0; i < 1000; i++ {
		h := RandomHash()
		data, _ := h.MarshalBinary()

		n := new(Hash)
		err := n.UnmarshalBinary(data)
		if err != nil {
			t.Error(err)
		}

		if !h.IsSameAs(n) {
			t.Error("Failed, should be same")
		}
	}

	h, err := HexToHash("")
	if err == nil {
		t.Error("Should fail")
	}

	h, err = HexToHash("0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Error(err)
	}

	if h.String() != "0000000000000000000000000000000000000000000000000000000000000000" {
		t.Error("Failed, should be all 0s")
	}
}
