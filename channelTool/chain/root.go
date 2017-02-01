package chain

import (
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

var _ = constants.SQL_DB

type RootChain struct {
	Register *factom.Entry
}

// Factom entry
//		byte		Version
//		[13]byte	Channel Chain"
//		[32]byte	rootChainID
func (r *RootChain) RegisterRootChain(rootChain primitives.Hash, publicKey3 primitives.PublicKey, sig []byte) {
	r.Register = new(factom.Entry)

	r.Register.ExtIDs = append(r.Register.ExtIDs, []byte{constants.FACTOM_VERSION})
	r.Register.ExtIDs = append(r.Register.ExtIDs, []byte("Channel Chain"))
	r.Register.ExtIDs = append(r.Register.ExtIDs, rootChain.Bytes())
	r.Register.ExtIDs = append(r.Register.ExtIDs, publicKey3.Bytes())
	r.Register.ExtIDs = append(r.Register.ExtIDs, sig)
}
