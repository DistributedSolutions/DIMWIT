package channelTool_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/channelTool"
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/FactomProject/factom"
)

var _ = fmt.Sprint("...")

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
	for i := 0; i < 10; i++ {
		c := common.RandomNewChannel()
		ec := factom.NewECAddress()
		a, err := NewAuthChannel(c, ec)
		if err != nil {
			t.Error(err)
		}

		chs, err := a.ReturnFactomChains()
		if err != nil {
			t.Error(err)
		}
		if len(chs) < 3 {
			t.Error("Should be at least 3 chains")
		}

		es, err := a.ReturnFactomEntries()
		if err != nil {
			t.Error(err)
		}
		if len(es) == 0 {
			t.Error("Should be more than 1 entry")
		}
	}
}
