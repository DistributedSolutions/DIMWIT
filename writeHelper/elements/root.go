package elements

import (
	//"fmt"
	// "time"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/writeHelper"
	"github.com/FactomProject/factom"
)

type Root struct {
	// Factom Elements
	Chain         *RootChain
	RegisterRoot  *RegisterRootEntry
	ContentSigKey *ContentKeyEntry
}

// InitiateChannel creates the factom elements needed to initiate a channel.
// It will set the channel variables and return the channel with new values
func InitiateChannel(a *writeHelper.AuthChannel, ch *common.Channel) (*writeHelper.AuthChannel, *common.Channel, *Root) {
	r := new(Root)
	r.Chain = new(RootChain)
	pubs := make([]primitives.PublicKey, 0)
	for _, p := range a.PrivateKeys {
		pubs = append(pubs, p.Public)
	}
	// RootChainID
	root := r.Chain.Create(pubs, ch.ChannelTitle)

	var _ = root
	return a, ch, r
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
func (rc *RootChain) Create(pubs []primitives.PublicKey, title primitives.Title) (root []byte) {
	rc.PubKeys = pubs
	rc.Title = title
	rc.Nonce, root = FindValidNonce(rc.AllButNonce())
	return
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

func (rre *RegisterRootEntry) Create(key3 primitives.PrivateKey, root primitives.Hash) {
	rre.RootChainID = root
	rre.KeyToSign = key3
}

func (rre *RegisterRootEntry) FactomEntry() *factom.Entry {
	e := new(factom.Entry)
	extIDs := VersionAndType(rre)
	extIDs = append(extIDs, rre.RootChainID.Bytes())
	extIDs = append(extIDs, rre.KeyToSign.Public.Bytes())
	sig := rre.KeyToSign.Sign(upToSig(extIDs))
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

func (cke *ContentKeyEntry) Create(key3 primitives.PrivateKey, root primitives.Hash, newConKey primitives.PublicKey) {
	cke.RootChainID = root
	cke.KeyToSign = key3
	cke.ContentPubKey = newConKey
}

func (cke *ContentKeyEntry) FactomEntry() *factom.Entry {
	e := new(factom.Entry)
	extIDs := VersionAndType(cke)
	extIDs = append(extIDs, cke.RootChainID.Bytes())
	extIDs = append(extIDs, cke.ContentPubKey.Bytes())
	extIDs = append(extIDs, TimeStampBytes())
	extIDs = append(extIDs, cke.KeyToSign.Public.Bytes())

	sig := cke.KeyToSign.Sign(upToSig(extIDs))
	extIDs = append(extIDs, sig)

	e.ExtIDs = extIDs
	e.Content = GetContentSignature()
	e.ChainID = cke.RootChainID.String()

	return e
}
