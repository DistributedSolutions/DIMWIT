package creation

import (
	"fmt"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

type ContentChain struct {
	FirstEntry *factom.Chain
	Entries    []*factom.Entry

	// Must be done after first entry
	Register *RegisterStruct
}

// Factom Chain
//		byte		Version
//		byte		ContentType
//		byte		TotalEntries
//		[13]byte	"Content Chain"
//		[32]byte	RootChainID
//		[20]byte	Infohash
//		[]byte		Timestamp
//		byte		ShiftCipher
//		[32]byte	ContentSignKey
//		[64]byte	Signature
//		[]byte		nonce
//	CONTENT TODO: CONTENT

//  CreateContentChain needs all the metadata to determine how many entries to use
func (r *ContentChain) CreateContentChain() {

}

// Factom Entry
//		byte		Version
//		byte		ContentType
//		[12]byte	"Content Link"
//		[32]byte	RootChainID
//		[]byte		Timestamp
//		[32]byte	ContentSignKey
//		[64]byte	Signature
func (r *ContentChain) RegisterNewContentChain(rootChain primitives.Hash, contentChainID primitives.Hash, contentType byte, sigKey primitives.PrivateKey) error {
	e := new(factom.Entry)

	timeData, err := time.Now().MarshalBinary()
	if err != nil {
		return fmt.Errorf("Unable to create a timestamp: %s", err.Error())
	}

	e.ExtIDs = append(e.ExtIDs, []byte{constants.FACTOM_VERSION}) // 0
	e.ExtIDs = append(e.ExtIDs, []byte{contentType})              // 1
	e.ExtIDs = append(e.ExtIDs, []byte("Content Link"))           // 2
	e.ExtIDs = append(e.ExtIDs, rootChain.Bytes())                // 3
	e.ExtIDs = append(e.ExtIDs, timeData)                         // 4
	e.ExtIDs = append(e.ExtIDs, sigKey.Public.Bytes())            // 5

	msg := upToNonce(e.ExtIDs, 4)
	sig := sigKey.Sign(msg)
	e.ExtIDs = append(e.ExtIDs, sig) // 4

	e.ChainID = contentChainID.String()
	r.Register.Entry = e

	return nil
}
