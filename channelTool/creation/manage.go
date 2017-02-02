package creation

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

type ManageChain struct {
	Register RegisterStruct
	Create   CreateStruct
}

// Factom Chain
//		byte		Version
//		[24]byte	"Channel Management Chain"
//		[32]byte	RootChainID
//		[]byte		Title
//		[]byte		nonce
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
func (r *ManageChain) CreateManagementChain(rootChain primitives.Hash, publicKeys []primitives.PublicKey, title primitives.Title, sigKey primitives.PrivateKey) error {
	r.endExtID = 4

	data, err := title.MarshalBinary()
	if err != nil {
		return err
	}

	e := new(factom.Entry)

	e.ExtIDs = append(e.ExtIDs, []byte{constants.FACTOM_VERSION})   // 0
	e.ExtIDs = append(e.ExtIDs, []byte("Channel Management Chain")) // 1
	e.ExtIDs = append(e.ExtIDs, rootChain.Bytes())                  // 2
	e.ExtIDs = append(e.ExtIDs, data)                               // 3
	r.Create.ExtIDs = e.ExtIDs
	nonce := FindValidNonce(r.Create)
	e.ExtIDs = append(e.ExtIDs, nonce)              // 4
	e.ExtIDs = append(e.ExtIDs, publicKeys.Bytes()) // 5

	msg := upToNonce(e.ExtIDs, 5)
	sig := sigKey.Sign(msg)
	e.ExtIDs = append(e.ExtIDs, sig) // 6

	r.Create.ExtIDs = e.ExtIDs

	c := factom.NewChain(e)
	r.Create.Chain = c
	return nil

}

// Factom Entry
//		byte		Version
//		[25]byte	"Register Management Chain"
//		[32]byte	ManagementChainID
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
func (r *RootChain) RegisterRootEntry(rootChain primitives.Hash, publicKey3 primitives.PublicKey, sig []byte) {
	e := new(factom.Entry)

	e.ExtIDs = append(e.ExtIDs, []byte{constants.FACTOM_VERSION})
	e.ExtIDs = append(e.ExtIDs, []byte("Register Management Chain"))
	e.ExtIDs = append(e.ExtIDs, rootChain.Bytes())
	e.ExtIDs = append(e.ExtIDs, publicKey3.Bytes())
	e.ExtIDs = append(e.ExtIDs, sig)

	r.Register.Entry = e
}
