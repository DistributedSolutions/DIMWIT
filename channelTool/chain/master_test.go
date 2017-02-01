package chain_test

import (
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/channelTool/chain"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
)

func TestMasterChain(t *testing.T) {
	m := NewMasterChain()
	//k := FindValidNonce(m)
	//fmt.Printf("%x", k)
	if m.Chain.ChainID != constants.MASTER_CHAIN_STRING {
		t.Error("Master chain ID does not match constants")
	}
}
