package objects_test

import (
	"testing"

	//"github.com/DistributedSolutions/DIMWIT/common"
	. "github.com/DistributedSolutions/DIMWIT/constructor/objects"
)

func TestChannelWrapper(t *testing.T) {
	for i := 0; i < 10; i++ {
		w := *RandomChannelWrapper()
		data, err := w.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		b := new(ChannelWrapper)
		nd, err := b.UnmarshalBinaryData(data)
		if err != nil {
			t.Error(err)
		}

		if !w.IsSameAs(b) {
			t.Error("Should be same")
		}

		if len(nd) != 0 {
			t.Errorf("Not unmarshaled correctly. Bytes %d bytes lefts over", len(nd))
		}
	}
}
