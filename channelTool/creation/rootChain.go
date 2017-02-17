package creation

import (
	"fmt"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

type RootChain struct {
	Register     RegisterStruct
	Create       CreateStruct
	ContSigEntry *factom.Entry
}

func (r *RootChain) ReturnChains() []*factom.Chain {
	c := make([]*factom.Chain, 0)
	c = append(c, r.Create.Chain)

	return c
}

func (r *RootChain) ReturnEntries() []*factom.Entry {
	c := make([]*factom.Entry, 0)
	c = append(c, r.Register.Entry)
	c = append(c, r.ContSigEntry)

	return c
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

	msg := upToNonce(e.ExtIDs)
	e.ExtIDs = append(e.ExtIDs, sigKey.Public.Bytes()) // 3
	sig := sigKey.Sign(msg)
	e.ExtIDs = append(e.ExtIDs, sig)
	e.ChainID = constants.MASTER_CHAIN_STRING

	r.Register.Entry = e
}

// Factom entry
//		byte		Version
//		[19]byte	"Content Signing Key"
//		[32]byte	RootChainID
//		[32]byte	ContentSigningKey
//		[15]byte	Timestamp
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
func (r *RootChain) ContentSigningKey(rootChain primitives.Hash, contentSigningKey primitives.PublicKey, sigKey primitives.PrivateKey) error {
	e := new(factom.Entry)

	ts := time.Now()
	tsData, err := ts.MarshalBinary()
	if err != nil {
		return err
	}

	e.ExtIDs = append(e.ExtIDs, []byte{constants.FACTOM_VERSION}) // 0
	e.ExtIDs = append(e.ExtIDs, []byte("Content Signing Key"))    // 1
	e.ExtIDs = append(e.ExtIDs, rootChain.Bytes())                // 2
	e.ExtIDs = append(e.ExtIDs, contentSigningKey.Bytes())        // 3
	e.ExtIDs = append(e.ExtIDs, tsData)                           // 4

	msg := upToNonce(e.ExtIDs)
	e.ExtIDs = append(e.ExtIDs, sigKey.Public.Bytes()) // 5
	sig := sigKey.Sign(msg)
	e.ExtIDs = append(e.ExtIDs, sig)
	e.ChainID = rootChain.String()

	r.ContSigEntry = e
	return nil
}
