package primitives_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
)

var _ = fmt.Sprintf("")

func TestPrivateKey(t *testing.T) {
	p, err := GeneratePrivateKey()
	if err != nil {
		t.Error(err)
	}

	sec := p.Secret[:32]
	p2 := new(PrivateKey)
	err = p2.SetBytes(sec)
	if err != nil {
		t.Error(err)
	}

	if !p2.Public.IsSameAs(&p.Public) {
		t.Error("Should be same")
	}

	var _, _ = p, sec
}

func TestPrivateKeyMarshal(t *testing.T) {
	for i := 0; i < 1000; i++ {
		h, _ := RandomPrivateKey()
		data, _ := h.MarshalBinary()

		n := new(PrivateKey)
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

		j := new(PrivateKey)
		err = json.Unmarshal(jdata, j)
		if err != nil {
			t.Error(err)
		}

		if !h.IsSameAs(j) {
			t.Errorf("Should be same. Found %s, expected %s", j.String(), h.String())
		}
	}
}

func TestCreateKey(t *testing.T) {
	// 4d801d9228505cbf1008b6a1a4d38edc3fcf40ecc08476d2c811813c0b239cfa
	// a2be234ba4bc77f58581e55d61d7db15018f188a0b31ce632b80918fc68dca13
	sec, _ := hex.DecodeString("4d801d9228505cbf1008b6a1a4d38edc3fcf40ecc08476d2c811813c0b239cfa")
	// pub, _ := hex.DecodeString("a2be234ba4bc77f58581e55d61d7db15018f188a0b31ce632b80918fc68dca13")

	pk, err := GeneratePrivateKeyFromHex("4d801d9228505cbf1008b6a1a4d38edc3fcf40ecc08476d2c811813c0b239cfa")
	if err != nil {
		t.Error(err)
	}

	if pk.String() != "4d801d9228505cbf1008b6a1a4d38edc3fcf40ecc08476d2c811813c0b239cfaa2be234ba4bc77f58581e55d61d7db15018f188a0b31ce632b80918fc68dca13" {
		t.Error("Private key failed to import from hex")
	}

	pk, err = GeneratePrivateKeyFromBytes(sec)
	if err != nil {
		t.Error(err)
	}

	if pk.String() != "4d801d9228505cbf1008b6a1a4d38edc3fcf40ecc08476d2c811813c0b239cfaa2be234ba4bc77f58581e55d61d7db15018f188a0b31ce632b80918fc68dca13" {
		t.Error("Private key failed to import from hex")
	}

	_, err = GeneratePrivateKeyFromBytes([]byte{})
	if err == nil {
		t.Error("should error")
	}

	_, err = GeneratePrivateKeyFromHex("aa")
	if err == nil {
		t.Error("should error")
	}

	_, err = GeneratePrivateKeyFromHex("Kd801d9228505cbf1008b6a1a4d38edc3fcf40ecc08476d2c811813c0b239cfa")
	if err == nil {
		t.Error("should error")
	}

}

func TestPrivateKeyDiff(t *testing.T) {
	same := 0
	for i := 0; i < 1000; i++ {
		a, _ := RandomPrivateKey()
		b, _ := RandomPrivateKey()
		if a.IsSameAs(b) {
			same++
		}
	}
	if same > 15 {
		t.Error("More than 15 are the same, it is totally random, so it is likely the IsSameAs is broken.")
	}
}

func TestEmptyPK(t *testing.T) {
	pk := new(PrivateKey)
	if !pk.Empty() {
		t.Error("Should be empty")
	}
}

func TestBadUnmarshalPrK(t *testing.T) {
	badData := []byte{}

	n := new(PrivateKey)

	err := n.UnmarshalBinary(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
