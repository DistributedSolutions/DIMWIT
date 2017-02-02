package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestMarshal(t *testing.T) {
	for i := 0; i < 1000; i++ {
		str := random.RandString()
		max := random.RandomIntBetween(0, 100)
		max += len(str)

		data, err := MarshalStringToBytes(str, max)
		if err != nil {
			t.Error(err)
		}

		resp, data, err := UnmarshalStringFromBytesData(data, max)
		if err != nil {
			t.Error(err)
		}

		if resp != str {
			t.Error("Unmarshal Fail")
		}

		if len(data) != 0 {
			t.Error("Unmarshal Return Data")
		}
	}
}

func TestUInt32Bytes(t *testing.T) {
	var i uint32 = 0
	for ; i < 2000; i++ {
		a := i
		data := Uint32ToBytes(a)

		b, err := BytesToUint32(data)
		if err != nil {
			t.Error(err)
		}

		if b != a {
			t.Error("Failed, should be same")
		}
	}

	//jesse testing for parts an series
	c := append([]byte{0x00, 0x00, 0x00}, []byte{0x01}...)
	data, err := BytesToUint32(c)
	if err != nil {
		t.Error(err)
	}
	if data != 1 {
		t.Errorf("1 != %d", data)
	}
	c = append([]byte{0x00, 0x00}, []byte{0x01, 0x00}...)
	data, err = BytesToUint32(c)
	if err != nil {
		t.Error(err)
	}
	if data != 256 {
		t.Errorf("256 != %d", data)
	}
}

func TestUInt64Bytes(t *testing.T) {
	var i uint64 = 0
	for ; i < 2000; i++ {
		a := i
		data, err := Uint64ToBytes(a)
		if err != nil {
			t.Error(err)
		}

		b, err := BytesToUint64(data)
		if err != nil {
			t.Error(err)
		}

		if b != a {
			t.Error("Failed, should be same")
		}
	}
}
