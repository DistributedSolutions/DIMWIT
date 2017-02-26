package primitives_test

import (
	"encoding/json"
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

		if i < 10 {
			err := n.UnmarshalBinary(data)
			if err != nil {
				t.Error(err)
			}
		}

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

		j := new(Hash)
		err = json.Unmarshal(jdata, j)
		if err != nil {
			t.Error(err)
		}

		if !h.IsSameAs(j) {
			t.Errorf("Should be same. Found %s, expected %s", j.String(), h.String())
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

	_, err = HexToHash("notvalid")
	if err == nil {
		t.Error("Not valid hex, should error")
	}

	_, err = HexToHash("000000000000000000000000000000000")
	if err == nil {
		t.Error("Not even length hex, should error")
	}

	err = i.UnmarshalBinary([]byte{0x00})
	if err == nil {
		t.Error("Should error")
	}

}

func TestBytesToHex(t *testing.T) {
	for i := 0; i < 1000; i++ {
		data := random.RandByteSliceOfSize(32)
		_, err := BytesToHash(data)
		if err != nil {
			t.Error(err)
		}
	}
	for i := 0; i < 1000; i++ {
		size := random.RandomIntBetween(0, 1000)
		if size == 32 {
			size += 1
		}
		data := random.RandByteSliceOfSize(size)
		_, err := BytesToHash(data)
		if err == nil {
			t.Error("Should Error")
		}
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

		if l.Empty() && len(l.GetHashes()) != 0 {
			t.Error("Should not be empty")
		}

		jdata, err := json.Marshal(n)
		if err != nil {
			t.Error(err)
		}

		j := new(HashList)
		err = json.Unmarshal(jdata, j)
		if err != nil {
			t.Error(err)
		}

		if !n.IsSameAs(j) {
			t.Errorf("Should be same. ")
		}
	}

	// Bad unmarshal
	h := new(HashList)
	err := h.UnmarshalBinary([]byte{0x00})
	if err == nil {
		t.Error("Should error")
	}
}

func TestDiffHashList(t *testing.T) {
	same := 0
	for i := 0; i < 1000; i++ {
		a := RandomHashList(random.RandomUInt32Between(0, 1000))
		b := RandomHashList(random.RandomUInt32Between(0, 1000))
		if a.IsSameAs(b) {
			same++
		}
	}
	if same > 15 {
		t.Error("More than 15 are the same, it is totally random, so it is likely the IsSameAs is broken.")
	}
}

func TestEmptyHash(t *testing.T) {
	hl := new(HashList)
	if !hl.Empty() {
		t.Error("Hashlist Should be empty")
	}

	h := new(Hash)
	if !h.Empty() {
		t.Error("Hash Should be empty")
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
