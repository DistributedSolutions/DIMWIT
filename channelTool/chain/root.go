package chain

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

var _ = constants.SQL_DB

type RootChain struct {
	Register RegisterStruct
	Create   CreateStruct
}

type RegisterStruct struct {
	Entry *factom.Entry
}

type CreateStruct struct {
	Chain  *factom.Chain
	ExtIDs [][]byte
}

// Factom entry
//		byte		Version
//		[13]byte	"Channel Chain"
//		[32]byte	RootChainID
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
func (r *RootChain) RegisterRootChain(rootChain primitives.Hash, publicKey3 primitives.PublicKey, sig []byte) {
	e := new(factom.Entry)

	e.ExtIDs = append(e.ExtIDs, []byte{constants.FACTOM_VERSION})
	e.ExtIDs = append(e.ExtIDs, []byte("Channel Chain"))
	e.ExtIDs = append(e.ExtIDs, rootChain.Bytes())
	e.ExtIDs = append(e.ExtIDs, publicKey3.Bytes())
	e.ExtIDs = append(e.ExtIDs, sig)

	r.Register.Entry = e
}

// Factom Chain
//		byte		Version
//		[18]byte	"Channel Root Chain"
//		[32]byte	PublicKey(1)
//		[32]byte	PublicKey(2)
//		[32]byte	PublicKey(3)
//		[]byte		Nonce
func (r *RootChain) CreateRootChain(publicKeys []primitives.PublicKey) error {
	e := new(factom.Entry)

	if len(publicKeys) != 3 {
		return fmt.Errorf("Not enough keys")
	}

	e.ExtIDs = append(e.ExtIDs, []byte{constants.FACTOM_VERSION}) // 0
	e.ExtIDs = append(e.ExtIDs, []byte("Channel Root Chain"))     // 1
	e.ExtIDs = append(e.ExtIDs, publicKeys[0].Bytes())            // 2
	e.ExtIDs = append(e.ExtIDs, publicKeys[1].Bytes())            // 3
	e.ExtIDs = append(e.ExtIDs, publicKeys[2].Bytes())            // 4
	r.Create.ExtIDs = e.ExtIDs
	nonce := FindValidNonce(r.Create)
	e.ExtIDs = append(e.ExtIDs, nonce) // 5
	r.Create.ExtIDs = e.ExtIDs

	c := factom.NewChain(e)
	r.Create.Chain = c
	return nil

}

func (c CreateStruct) upToNonce() []byte {
	return upToNonce(c.ExtIDs, 5)
}
