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
