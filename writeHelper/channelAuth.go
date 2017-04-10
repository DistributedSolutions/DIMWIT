package writeHelper

import (
	//"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common"
	//"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	//"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
	"github.com/DistributedSolutions/DIMWIT/writeHelper/elements"
	"github.com/FactomProject/factom"
)

const PRIV_KEY_AMT int = 3

type AuthChannel struct {
	ChannelRoot   primitives.Hash
	ChannelManage primitives.Hash

	PrivateKeys    [PRIV_KEY_AMT]primitives.PrivateKey
	ContentSigning primitives.PrivateKey

	EntryCreditKey *factom.ECAddress

	// Non-Marshaled
	Root   *elements.Root
	Manage *elements.Manage
}

func NewAuthChannel(ch *common.Channel, ec *factom.ECAddress) (*AuthChannel, error) {
	if !factom.IsValidAddress(ec.String()) {
		return nil, fmt.Errorf("Entry credit address is invalid")
	}
	var err error

	if ec == nil {
		return nil, fmt.Errorf("ECAddress is nil")
	}

	a := new(AuthChannel)

	// TODO: Discover Root
	//a.ChannelRoot = root

	for i := 0; i < PRIV_KEY_AMT; i++ {
		pk, err := primitives.GeneratePrivateKey()
		if err != nil {
			return nil, err
		}
		a.PrivateKeys[i] = *pk
		switch i {
		case 0:
			ch.LV1PublicKey = pk.Public
		case 1:
			ch.LV2PublicKey = pk.Public
		case 2:
			ch.LV3PublicKey = pk.Public
		}
	}

	pk, err := primitives.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	a.ContentSigning = *pk
	ch.ContentSingingKey = pk.Public
	a.EntryCreditKey = ec

	return a, nil
}

func (a *AuthChannel) Initiate(ch *common.Channel) {
	r := a.initRoot(ch)
	m := a.initManage(ch)

	a.Root = r
	a.Manage = m
	return
}

func (a *AuthChannel) initRoot(ch *common.Channel) *elements.Root {
	r := elements.NewRoot()
	pubs := make([]primitives.PublicKey, 0)
	for _, p := range a.PrivateKeys {
		pubs = append(pubs, p.Public)
	}

	// RootChainID
	root, err := r.Chain.Create(pubs, ch.ChannelTitle)
	if err != nil { // TODO: Handle error
		return nil
	}

	r.RegisterRoot.Create(a.PrivateKeys[2], root)
	r.ContentSigKey.Create(a.PrivateKeys[2], root, a.ContentSigning.Public)

	a.ChannelRoot = *root
	ch.RootChainID = *root
	return r
}

func (a *AuthChannel) initManage(ch *common.Channel) *elements.Manage {
	m := elements.NewManage()
	man, err := m.ManageChain.Create(a.ChannelRoot, a.PrivateKeys[2])
	if err != nil { // TODO: Handle error
		return nil
	}

	m.RegisterManageChain.Create(a.ChannelRoot, *man, a.PrivateKeys[2])

	a.ChannelManage = *man
	ch.ManagementChainID = *man
	return m
}
