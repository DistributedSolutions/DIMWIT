package creation

// Master chain. This is only called once, ever, but nice to use in testing.
// Also a nice blank reference

import (
	//"bytes"
	//"crypto/sha256"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/FactomProject/factom"
)

type MasterChain struct {
	Chain *factom.Chain
}

func NewMasterChain() *MasterChain {
	m := new(MasterChain)

	e := new(factom.Entry)
	e.ExtIDs = append(e.ExtIDs, []byte{constants.FACTOM_VERSION})
	e.ExtIDs = append(e.ExtIDs, []byte("Master Chain"))
	// Hex: 313730313439
	// Int: 170149
	e.ExtIDs = append(e.ExtIDs, []byte{0x31, 0x37, 0x30, 0x31, 0x34, 0x39})

	c := factom.NewChain(e)
	m.Chain = c

	return m
}

/*
func (m *MasterChain) upToNonce() []byte {
	buf := new(bytes.Buffer)

	result := sha256.Sum256([]byte{constants.FACTOM_VERSION})
	buf.Write(result[:])

	result = sha256.Sum256([]byte("Master Chain"))
	buf.Write(result[:])

	return buf.Next(buf.Len())
}

func (m *MasterChain) getNonce() []byte {
	return m.Chain.FirstEntry.ExtIDs[2]
}
*/
