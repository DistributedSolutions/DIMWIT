package creation_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/channelTool/creation"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
)

var _ = fmt.Sprintf("")

func TestRoot(t *testing.T) {
	p := make([]primitives.PublicKey, 3)
	for i := range p {
		p[i] = *primitives.RandomPublicKey()
	}

	rc := new(RootChain)
	rc.CreateRootChain(p)

	data, err := hex.DecodeString(rc.Create.Chain.FirstEntry.ChainID)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if bytes.Compare(data[:constants.CHAIN_PREFIX_LENGTH_CHECK], constants.CHAIN_PREFIX[:constants.CHAIN_PREFIX_LENGTH_CHECK]) != 0 {
		t.Error("Invalid chainID")
	}
}
