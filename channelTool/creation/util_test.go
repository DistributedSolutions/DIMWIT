package creation_test

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/channelTool/creation"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

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
