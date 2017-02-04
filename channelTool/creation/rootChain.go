package creation

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

type RootChain struct {
	Register RegisterStruct
	Create   CreateStruct
}

// Factom Chain
//		byte		Version
//		[18]byte	"Channel Root Chain"
//		[]byte		Title
//		[32]byte	PublicKey(1)
//		[32]byte	PublicKey(2)
//		[32]byte	PublicKey(3)
//		[]byte		Nonce
func (r *RootChain) CreateRootChain(publicKeys []primitives.PublicKey, title primitives.Title) error {
	r.Create.endExtID = 6

	titleData, err := title.MarshalBinary()
	if err != nil {
		return err
	}

	e := new(factom.Entry)

	if len(publicKeys) != 3 {
		return fmt.Errorf("Not enough keys")
	}

	e.ExtIDs = append(e.ExtIDs, []byte{constants.FACTOM_VERSION}) // 0
	e.ExtIDs = append(e.ExtIDs, []byte("Channel Root Chain"))     // 1
	e.ExtIDs = append(e.ExtIDs, titleData)                        // 2
	e.ExtIDs = append(e.ExtIDs, publicKeys[0].Bytes())            // 3
	e.ExtIDs = append(e.ExtIDs, publicKeys[1].Bytes())            // 4
	e.ExtIDs = append(e.ExtIDs, publicKeys[2].Bytes())            // 5
	r.Create.ExtIDs = e.ExtIDs
	nonce := FindValidNonce(r.Create)
	e.ExtIDs = append(e.ExtIDs, nonce) // 6
	r.Create.ExtIDs = e.ExtIDs

	c := factom.NewChain(e)
	r.Create.Chain = c
	return nil

}

// Factom entry
//		byte		Version
//		[13]byte	"Channel Chain"
//		[32]byte	RootChainID
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
func (r *RootChain) RegisterRootEntry(rootChain primitives.Hash, sigKey primitives.PrivateKey) {
	e := new(factom.Entry)

	e.ExtIDs = append(e.ExtIDs, []byte{constants.FACTOM_VERSION}) // 0
	e.ExtIDs = append(e.ExtIDs, []byte("Channel Chain"))          // 1
	e.ExtIDs = append(e.ExtIDs, rootChain.Bytes())                // 2
	e.ExtIDs = append(e.ExtIDs, sigKey.Public.Bytes())            // 3

	msg := upToNonce(e.ExtIDs, 4)
	sig := sigKey.Sign(msg)
	e.ExtIDs = append(e.ExtIDs, sig)

	r.Register.Entry = e
}
