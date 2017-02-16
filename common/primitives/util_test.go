package primitives_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = fmt.Sprintf("")

func TestPropString(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("Testing Marshing string", prop.ForAll(
		func(t1 string) bool {
			data, err1 := MarshalStringToBytes(t1, len(t1)+1)
			t2, newData, err2 := UnmarshalStringFromBytesData(data, len(t1)+1)
			return err1 == nil && err2 == nil && len(newData) == 0 && t1 == t2
		},
		gen.AnyString(),
	))

	//properties.Run(gopter.ConsoleReporter(true))
	properties.TestingRun(t)
}

func TestMarshal(t *testing.T) {
	for i := 0; i < 1000; i++ {
		str := random.RandString()
		if i > 950 {
			str = string(0x00) + str
			str = str + string(0x00)
		}
		max := random.RandomIntBetween(0, 100)
		max += len(str)

		data, err := MarshalStringToBytes(str, max)
		if err != nil {
			t.Error(err)
		}

		if i < 10 {
			_, err := UnmarshalStringFromBytes(data, max)
			if err != nil {
				t.Error(err)
			}
		}

		resp, data, err := UnmarshalStringFromBytesData(data, max)
		if err != nil {
			t.Error(err)
		}

		str = strings.Replace(str, string(0x00), string(0x01), -1)
		if resp != str {
			t.Error("Unmarshal Fail")
		}

		if len(data) != 0 {
			t.Error("Unmarshal Return Data")
		}
	}

	str := "123456"

	data, err := MarshalStringToBytes(str, 2)
	if err == nil {
		t.Error("Should error")
	}

	data, err = MarshalStringToBytes(str, 10)
	if err != nil {
		t.Error(err)
	}

	_, _, err = UnmarshalStringFromBytesData(data, 2)
	if err == nil {
		t.Error("should error")
	}

	// Bad marshal
	str, err = UnmarshalStringFromBytes([]byte{}, 0)
	if err == nil {
		t.Error("Should error")
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

	// Bad Marshal
	_, err := BytesToUint64([]byte{0x00, 0x00})
	if err == nil {
		t.Error("Should error")
	}
}

func TestXORCipher(t *testing.T) {
	for i := 0; i < 10000; i++ {
		var msg []byte
		for len(msg) == 0 {
			msg = random.RandByteSlice()
		}
		key := msg[0]
		cryp := XORCipher(key, msg)

		err := checkXOREnc(key, msg, cryp)
		if err != nil {
			t.Error(err)
		}

		decode := XORCipher(key, cryp)
		if bytes.Compare(decode, msg) != 0 {
			t.Error("Decode failed. Bytes are different")
		}
	}

}

func checkXOREnc(key byte, msg []byte, cryp []byte) error {
	for i := range msg {
		switch {
		case key == 0x00: // Should be same
			if msg[i] != cryp[i] {
				return fmt.Errorf("XORCipher failed. Key is 0x00, but msg byte is %x, and crypt byte is %x. Crypt should be %x", msg[i], cryp[i], msg[i])
			}
		case key == msg[i]:
			if cryp[i] != 0x00 {
				return fmt.Errorf("XORCipher failed. Key and msg byte is %x, and crypt byte is %x. Crypt should be 0x00", msg[i], cryp[i])
			}
		}
	}
	return nil
}

func TestRandXORKey(t *testing.T) {
	for i := 0; i < 100000; i++ {
		k := RandXORKey()
		if k == 0x00 {
			t.Error("Rand key should never be 0x00")
		}
	}
}
