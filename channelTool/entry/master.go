package entry

// Master chain. This is only called once, ever, but nice to use in testing.
// Also a nice blank reference

import (
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/FactomProject/factom"
)

type MasterChain struct {
	Entry factom.Entry
}

func NewMasterChain() *MasterChain {
	m := new(MasterChain)
	m.Entry.ExtIDs = append(m.Entry.ExtIDs, []byte{constants.FACTOM_VERSION})
	m.Entry.ExtIDs = append(m.Entry.ExtIDs, []byte("Master Chain"))
	m.Entry.ExtIDs = append(m.Entry.ExtIDs, []byte{0x01, 0x02})

	return m
}
