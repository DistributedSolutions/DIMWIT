package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/common/primitives"
)

var _ = fmt.Sprintf("")

func TestInfoHash(t *testing.T) {
	i, err := HexToHash("")
	if err == nil {
		t.Error("Should fail")
	}

	i, err = HexToHash("0000000000000000000000000000000000000000")
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%x\n", i.Bytes())
}
