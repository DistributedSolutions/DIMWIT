package primitives_test

import (
	"encoding/json"
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

		if i > 10 {
			continue
		}

		jdata, err := json.Marshal(h)
		if err != nil {
			t.Error(err)
		}

		j := new(PublicKey)
		err = json.Unmarshal(jdata, j)
		if err != nil {
			t.Error(err)
		}

		if !h.IsSameAs(j) {
			t.Errorf("Should be same. Found %s, expected %s", j.String(), h.String())
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

func TestPublicKeyDiff(t *testing.T) {
	same := 0
	for i := 0; i < 1000; i++ {
		a := RandomPublicKey()
		b := RandomPublicKey()
		if a.IsSameAs(b) {
			same++
		}
	}
	if same > 15 {
		t.Error("More than 15 are the same, it is totally random, so it is likely the IsSameAs is broken.")
	}
}

func TestEmptyPub(t *testing.T) {
	pk := new(PublicKey)
	if !pk.Empty() {
		t.Error("Should be empty")
	}
}

func TestBadUnmarshalPK(t *testing.T) {
	badData := []byte{}

	n := new(PublicKey)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
