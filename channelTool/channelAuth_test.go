package channelTool_test

import (
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/channelTool"
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/FactomProject/factom"
)

func TestAuthChannel(t *testing.T) {
	for i := 0; i < 100; i++ {
		a, err := RandomAuthChannel()
		if err != nil {
			t.Error(err)
			continue
		}

		data, err := a.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		b := new(AuthChannel)
		newData, err := b.UnmarshalBinaryData(data)
		if err != nil {
			t.Error(err)
		} else if len(newData) != 0 {
			t.Error("Should be no bytes left")
		}

		if !b.IsSameAs(a) {
			t.Error("Should be same")
		}
	}
}

func TestCompleteAuthChannel(t *testing.T) {
	for i := 0; i < 1000; i++ {
		c := common.RandomNewChannel()
		ec := factom.NewECAddress()

		_, err := NewAuthChannel(c, ec)
		if err != nil {
			t.Error(err)
		}
	}
}
