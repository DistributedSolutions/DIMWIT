package chain_test

import (
	"fmt"
	"testing"

	. "github.com/DistributedSolutions/DIMWIT/channelTool/chain"
)

func TestMasterChain(t *testing.T) {
	m := NewMasterChain()
	//k := FindValidNonce(m)
	//fmt.Printf("%x", k)
	fmt.Printf("%s\n", m.Chain.ChainID)
}
