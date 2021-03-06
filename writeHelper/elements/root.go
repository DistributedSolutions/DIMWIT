package elements

import (
	"fmt"
	// "time"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

var _ = fmt.Printf

type Root struct {
	Chain         *RootChain
	RegisterRoot  *RegisterRootEntry
	ContentSigKey *ContentKeyEntry
}

func NewRoot() *Root {
	r := new(Root)
	r.Chain = new(RootChain)
	r.RegisterRoot = new(RegisterRootEntry)
	r.ContentSigKey = new(ContentKeyEntry)

	return r
}

func (r *Root) FactomElements() ([]*factom.Entry, *factom.Chain) {
	es := make([]*factom.Entry, 0)

	es = append(es, r.RegisterRoot.FactomEntry())
	es = append(es, r.ContentSigKey.FactomEntry())

	return es, r.Chain.FactomChain()
}

// Factom Chain
//		byte		Version
//		[18]byte	"Channel Root Chain"
//		[]byte		Title
//		[32]byte	PublicKey(1)
//		[32]byte	PublicKey(2)
//		[32]byte	PublicKey(3)
//		[]byte		Nonce
type RootChain struct {
	Title   primitives.Title
	PubKeys []primitives.PublicKey
	Nonce   []byte
}

func (RootChain) Type() []byte  { return TYPE_ROOT_CHAIN }
func (RootChain) IsChain() bool { return true }
func (RootChain) ForChain() int { return CHAIN_NA }

// Create will find the nonce and return the root chain ID
func (rc *RootChain) Create(pubs []primitives.PublicKey, title primitives.Title) (rootHash *primitives.Hash, err error) {
	var root []byte
	rc.PubKeys = pubs
	rc.Title = title
	rc.Nonce, root = FindValidNonce(rc.AllButNonce())
	return primitives.BytesToHash(root)
}

func (rc *RootChain) AllButNonce() [][]byte {
	extIDs := VersionAndType(rc)
	data, _ := rc.Title.MarshalBinary()
	extIDs = append(extIDs, data)
	for _, p := range rc.PubKeys {
		extIDs = append(extIDs, p.Bytes())
	}

	return extIDs
}

func (rc *RootChain) FactomChain() *factom.Chain {
	e := new(factom.Entry)
	e.ExtIDs = append(rc.AllButNonce(), rc.Nonce)
	e.Content = GetContentSignature()
	return factom.NewChain(e)
}

// Factom entry
//		byte		Version
//		[13]byte	"Channel Chain"
//		[32]byte	RootChainID
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
type RegisterRootEntry struct {
	RootChainID primitives.Hash
	KeyToSign   primitives.PrivateKey
}

func (RegisterRootEntry) Type() []byte  { return TYPE_ROOT_REGISTER }
func (RegisterRootEntry) IsChain() bool { return false }
func (RegisterRootEntry) ForChain() int { return CHAIN_MAIN }

func (rre *RegisterRootEntry) Create(key3 primitives.PrivateKey, root *primitives.Hash) {
	rre.RootChainID = *root
	rre.KeyToSign = key3
}

func (rre *RegisterRootEntry) FactomEntry() *factom.Entry {
	e := new(factom.Entry)
	extIDs := VersionAndType(rre)
	extIDs = append(extIDs, rre.RootChainID.Bytes())

	sig := rre.KeyToSign.Sign(upToSig(extIDs))
	extIDs = append(extIDs, rre.KeyToSign.Public.Bytes())
	extIDs = append(extIDs, sig)

	e.ExtIDs = extIDs
	e.Content = GetContentSignature()
	e.ChainID = constants.MASTER_CHAIN_STRING

	return e
}

// Factom entry
//		byte		Version
//		[19]byte	"Content Signing Key"
//		[32]byte	RootChainID
//		[32]byte	ContentSigningKey
//		[15]byte	Timestamp
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
type ContentKeyEntry struct {
	RootChainID   primitives.Hash
	KeyToSign     primitives.PrivateKey
	ContentPubKey primitives.PublicKey
}

func (ContentKeyEntry) Type() []byte  { return TYPE_ROOT_CONTENT_KEY }
func (ContentKeyEntry) IsChain() bool { return false }
func (ContentKeyEntry) ForChain() int { return CHAIN_ROOT }

func (cke *ContentKeyEntry) Create(key3 primitives.PrivateKey, root *primitives.Hash, newConKey primitives.PublicKey) {
	cke.RootChainID = *root
	cke.KeyToSign = key3
	cke.ContentPubKey = newConKey
}

func (cke *ContentKeyEntry) FactomEntry() *factom.Entry {
	e := new(factom.Entry)
	extIDs := VersionAndType(cke)
	extIDs = append(extIDs, cke.RootChainID.Bytes())
	extIDs = append(extIDs, cke.ContentPubKey.Bytes())
	extIDs = append(extIDs, TimeStampBytes())
	sig := cke.KeyToSign.Sign(upToSig(extIDs))

	extIDs = append(extIDs, cke.KeyToSign.Public.Bytes())
	extIDs = append(extIDs, sig)

	e.ExtIDs = extIDs
	e.Content = GetContentSignature()
	e.ChainID = cke.RootChainID.String()

	return e
}
