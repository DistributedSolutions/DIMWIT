package elements

import (
	//"bytes"
	//"crypto/sha256"
	//"fmt"
	//"time"

	//	"github.com/DistributedSolutions/DIMWIT/common"
	//	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

type ChanContentChain struct {
	ContentChain          *CContentChain
	RegisterCContentChain *RegisterCContentChain
}

func NewChanContentChain() *ChanContentChain {
	m := new(ChanContentChain)
	m.ContentChain = new(CContentChain)
	m.RegisterCContentChain = new(RegisterCContentChain)

	return m
}

func (ccc *ChanContentChain) FactomElements() ([]*factom.Entry, *factom.Chain) {
	es := make([]*factom.Entry, 0)
	c := ccc.ContentChain.FactomChain()

	es = append(es, ccc.RegisterCContentChain.FactomEntry())

	return es, c
}

// Factom Chain
//		byte		Version
//		[24]byte	"Channel Content Chain"
//		[32]byte	RootChainID
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
//		[]byte		nonce
type CContentChain struct {
	RootChainID primitives.Hash
	KeyToSign   primitives.PrivateKey
	Nonce       []byte
}

func (m *CContentChain) Create(root primitives.Hash, key3 primitives.PrivateKey) (*primitives.Hash, error) {
	m.RootChainID = root

	nonce, chainID := FindValidNonce(m.AllButNonce())
	m.Nonce = nonce

	return primitives.BytesToHash(chainID)
}

func (m *CContentChain) AllButNonce() [][]byte {
	extIDs := VersionAndType(m)
	extIDs = append(extIDs, m.RootChainID.Bytes())

	sig := m.KeyToSign.Sign(upToSig(extIDs))
	extIDs = append(extIDs, m.KeyToSign.Public.Bytes())
	extIDs = append(extIDs, sig)

	return extIDs
}

func (m *CContentChain) FactomChain() *factom.Chain {
	e := new(factom.Entry)
	extIDs := m.AllButNonce()
	e.ExtIDs = append(extIDs, m.Nonce)

	return factom.NewChain(e)
}

func (CContentChain) Type() []byte  { return TYPE_CHANNEL_CONENT_CHAIN }
func (CContentChain) IsChain() bool { return true }
func (CContentChain) ForChain() int { return CHAIN_NA }

// Factom Entry
//		byte		Version
//		[25]byte	"Register Content Chain"
//		[32]byte	Channel Content ChainID
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
type RegisterCContentChain struct {
	RootChainID     primitives.Hash
	CContentChainID primitives.Hash
	KeyToSign       primitives.PrivateKey
}

func (RegisterCContentChain) Type() []byte  { return TYPE_CHANNEL_CONENT_CHAIN_REGISTER }
func (RegisterCContentChain) IsChain() bool { return false }
func (RegisterCContentChain) ForChain() int { return CHAIN_ROOT }

func (rmc *RegisterCContentChain) Create(rootChain primitives.Hash, contentChanID primitives.Hash, key3 primitives.PrivateKey) {
	rmc.RootChainID = rootChain
	rmc.CContentChainID = contentChanID
	rmc.KeyToSign = key3
}

func (rmc *RegisterCContentChain) FactomEntry() *factom.Entry {
	extIDs := VersionAndType(rmc)
	extIDs = append(extIDs, rmc.CContentChainID.Bytes())

	sig := rmc.KeyToSign.Sign(upToSig(extIDs))
	extIDs = append(extIDs, rmc.KeyToSign.Public.Bytes())
	extIDs = append(extIDs, sig)

	e := new(factom.Entry)
	e.ExtIDs = extIDs
	e.ChainID = rmc.RootChainID.String()
	return e
}
