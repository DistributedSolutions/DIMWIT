package random_test

import (
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

func TestRandInts(t *testing.T) {
	for i := 0; i < 10000; i++ {
		n := RandomIntBetween(0, 1000)
		if n < 0 || n > 1000 {
			t.Errorf("RandIntBetween failed. Should be between 0 and 1000, found %d", n)
		}

		bn := RandomInt64Between(666, 10000)
		if bn < 666 || bn > 10000 {
			t.Errorf("RandIntBetween failed. Should be between 666 and 10000, found %d", bn)
		}

		un := RandomUInt32Between(2301, 23455)
		if un < 2301 || un > 23455 {
			t.Errorf("RandIntBetween failed. Should be between 2301 and 23455, found %d", bn)
		}

		var _ = RandomUInt32()
	}
}

func TestRandString(t *testing.T) {
	for i := 0; i < 10000; i++ {
		l := RandomIntBetween(0, 500)
		s := RandStringOfSize(l)
		if len(s) != l {
			t.Errorf("RandStringOfSize Failed. String size is %d, should be %d", len(s), l)
		}

		h, err := RandomHexStringOfSize(l)
		if l%2 == 1 {
			if err == nil {
				t.Errorf("Should error as hex must be even length. Length given is %d", l)
			}
		} else {
			if len(h) != l {
				t.Errorf("RandomHexStringOfSize Failed. String size is %d, should be %d", len(h), l)
			}
		}
	}
}

func TestRandBytes(t *testing.T) {
	for i := 0; i < 10000; i++ {
		l := RandomIntBetween(0, 500)
		data := RandByteSliceOfSize(l)
		if len(data) != l {
			t.Errorf("RandByteSliceOfSize Failed. ByteSlice size is %d, should be %d", len(data), l)
		}
	}
}
