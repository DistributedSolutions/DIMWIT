package primitives_test

import (
	"encoding/json"
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

		jdata, err := json.Marshal(h)
		if err != nil {
			t.Error(err)
		}

		j := new(InfoHash)
		err = json.Unmarshal(jdata, j)
		if err != nil {
			t.Error(err)
		}

		if !h.IsSameAs(j) {
			t.Errorf("Should be same. Found %s, expected %s", j.String(), h.String())
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

func TestInfoHashDiff(t *testing.T) {
	same := 0
	for i := 0; i < 1000; i++ {
		a := RandomInfoHash()
		b := RandomInfoHash()
		if a.IsSameAs(b) {
			same++
		}
	}
	if same > 15 {
		t.Error("More than 15 are the same, it is totally random, so it is likely the IsSameAs is broken.")
	}
}

func TestEmptyIH(t *testing.T) {
	ih := new(InfoHash)
	if !ih.Empty() {
		t.Error("Should be empty")
	}
}

func TestBadUnmarshalInfoHash(t *testing.T) {
	badData := []byte{}

	n := new(InfoHash)

	err := n.UnmarshalBinary(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}

	err = n.UnmarshalJSON([]byte{0x00})
	if err == nil {
		t.Error("Should panic or error out")
	}

	err = n.UnmarshalJSON([]byte("T28653dd34894ecff4fbbaf2dc513aae917659a2"))
	if err == nil {
		t.Error("Should panic or error out")
	}

	_, err = HexToInfoHash("T28653dd34894ecff4fbbaf2dc513aae917659a2")
	if err == nil {
		t.Error("Should panic or error out")
	}
}
