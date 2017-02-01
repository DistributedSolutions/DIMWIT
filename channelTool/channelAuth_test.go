package channelTool_test

import (
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/channelTool"
)

func TestAuthChannel(t *testing.T) {
	for i := 0; i < 1; i++ {
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
