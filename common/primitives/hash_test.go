package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
)

var _ = fmt.Sprintf("")

func TestHash(t *testing.T) {
	h, err := HexToHash("")
	if err == nil {
		t.Error("Should fail")
	}

	h, err = HexToHash("0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Error(err)
	}

	if h.String() != "0000000000000000000000000000000000000000000000000000000000000000" {
		t.Error("Failed, should be all 0s")
	}
}
