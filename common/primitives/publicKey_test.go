package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestPublicKey(t *testing.T) {
	for i := 0; i < 1000; i++ {
		h := RandomPublicKey()
		data, _ := h.MarshalBinary()

		n := new(PublicKey)
		err := n.UnmarshalBinary(data)
		if err != nil {
			t.Error(err)
		}

		if !h.IsSameAs(n) {
			t.Error("Failed, should be same")
		}
	}

	m := new(PublicKey)

	i, err := PublicKeyFromHex("")
	if err == nil {
		t.Error("Should fail")
	}

	str, _ := random.RandomHexStringOfSize(m.Length() * 2)
	i, err = PublicKeyFromHex(str)
	if err != nil {
		t.Error(err)
	}

	if i.String() != str {
		t.Error("Failed, should be all 0s")
	}
}
