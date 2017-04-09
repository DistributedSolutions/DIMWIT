package writeHelper

import (
	//"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common"
	//"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	//"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
	"github.com/FactomProject/factom"
)

const PRIV_KEY_AMT int = 3

type AuthChannel struct {
	ChannelRoot primitives.Hash

	PrivateKeys    [PRIV_KEY_AMT]primitives.PrivateKey
	ContentSigning primitives.PrivateKey

	EntryCreditKey *factom.ECAddress
}

func NewAuthChannel(ch *common.Channel, ec *factom.ECAddress) (*AuthChannel, *common.Channel, error) {
	if !factom.IsValidAddress(ec.String()) {
		return nil, nil, fmt.Errorf("Entry credit address is invalid")
	}
	var err error

	if ec == nil {
		return nil, nil, fmt.Errorf("ECAddress is nil")
	}

	a := new(AuthChannel)

	// TODO: Discover Root
	//a.ChannelRoot = root

	for i := 0; i < PRIV_KEY_AMT; i++ {
		pk, err := primitives.GeneratePrivateKey()
		if err != nil {
			return nil, nil, err
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
		return nil, nil, err
	}

	a.ContentSigning = *pk
	ch.ContentSingingKey = pk.Public
	a.EntryCreditKey = ec

	return a, ch, nil
}
