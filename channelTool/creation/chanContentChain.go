package creation

import (
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

type ChanContentChain struct {
	Register RegisterStruct
	Create   CreateStruct
}

// Factom Chain
//		byte		Version
//		[24]byte	"Channel Content Chain"
//		[32]byte	RootChainID
//		[]byte		Title
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
//		[]byte		nonce
func (r *ChanContentChain) CreateContentChain(rootChain primitives.Hash, sigKey primitives.PrivateKey) error {
	r.Create.endExtID = 5

	e := new(factom.Entry)

	e.ExtIDs = append(e.ExtIDs, []byte{constants.FACTOM_VERSION})   // 0
	e.ExtIDs = append(e.ExtIDs, []byte("Channel Management Chain")) // 1
	e.ExtIDs = append(e.ExtIDs, rootChain.Bytes())                  // 2
	e.ExtIDs = append(e.ExtIDs, sigKey.Public.Bytes())              // 3

	msg := upToNonce(e.ExtIDs, 4)
	sig := sigKey.Sign(msg)
	e.ExtIDs = append(e.ExtIDs, sig) // 4

	r.Create.ExtIDs = e.ExtIDs
	nonce := FindValidNonce(r.Create)
	e.ExtIDs = append(e.ExtIDs, nonce) // 5
	r.Create.ExtIDs = e.ExtIDs

	c := factom.NewChain(e)
	r.Create.Chain = c
	return nil

}

// Factom Entry
//		byte		Version
//		[25]byte	"Register Content Chain"
//		[32]byte	Channel Content ChainID
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
func (r *ChanContentChain) RegisterChannelContentChain(rootChain primitives.Hash, contentChainID primitives.Hash, sigKey primitives.PrivateKey) {
	e := new(factom.Entry)

	e.ExtIDs = append(e.ExtIDs, []byte{constants.FACTOM_VERSION}) // 0
	e.ExtIDs = append(e.ExtIDs, []byte("Register Content Chain")) // 1
	e.ExtIDs = append(e.ExtIDs, contentChainID.Bytes())           // 2
	e.ExtIDs = append(e.ExtIDs, sigKey.Public.Bytes())            // 3

	msg := upToNonce(e.ExtIDs, 3)
	sig := sigKey.Sign(msg)
	e.ExtIDs = append(e.ExtIDs, sig) // 4

	e.ChainID = rootChain.String()
	r.Register.Entry = e
}