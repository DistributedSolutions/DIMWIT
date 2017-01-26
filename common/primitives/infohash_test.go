package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
)

var _ = fmt.Sprintf("")

func TestInfoHash(t *testing.T) {
	i, err := HexToInfoHash("")
	if err == nil {
		t.Error("Should fail")
	}

	i, err = HexToInfoHash("0000000000000000000000000000000000000000")
	if err != nil {
		t.Error(err)
	}

	if i.String() != "0000000000000000000000000000000000000000" {
		t.Error("Failed, should be all 0s")
	}
}
