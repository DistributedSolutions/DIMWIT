package channelTool

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/channelTool/creation"
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

const PRIV_KEY_AMT int = 3

type AuthChannel struct {
	Channel common.Channel

	PrivateKeys    [PRIV_KEY_AMT]primitives.PrivateKey
	ContentSigning primitives.PrivateKey

	EntryCreditKey *factom.ECAddress

	// Not marshaled, timestamps in these.
	// Have all factom entries
	RootChain    *creation.RootChain
	ManageChain  *creation.ManageChain
	ContentChain *creation.ChanContentChain
}

// Makes the authority channel and builds all factom components
func NewAuthChannel(ch *common.Channel, ec *factom.ECAddress) (*AuthChannel, error) {
	if ch.Status() < constants.CHANNEL_READY {
		return nil, fmt.Errorf("Channel given is not ready, it is missing elements")
	}

	if !factom.IsValidAddress(ec.String()) {
		return nil, fmt.Errorf("Entry credit address is invalid")
	}

	a := new(AuthChannel)
	a.Channel = *ch

	a.PrivateKeys = make([]primitives.PrivateKey, PRIV_KEY_AMT)
	for i := 0; i < PRIV_KEY_AMT; i++ {
		pk, err := primitives.GeneratePrivateKey()
		if err != nil {
			return nil, err
		}
		a.PrivateKeys[i] = *pk
	}

	pk, err := primitives.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	a.ContentSigning = *pk
	a.EntryCreditKey = ec

	err = a.MakeChannel()
	if err != nil {
		return nil, err
	}

	err = a.MakeManagerChain()
	if err != nil {
		return nil, err
	}

	err = a.MakeContentChain()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func RandomAuthChannel() (*AuthChannel, error) {
	c := new(AuthChannel)
	var err error
	c.Channel = *common.RandomNewChannel()

	for i := 0; i < PRIV_KEY_AMT; i++ {
		sec, err := primitives.RandomPrivateKey()
		if err != nil {
			return nil, err
		}
		c.PrivateKeys[i] = *sec
	}

	sec, err := primitives.RandomPrivateKey()
	if err != nil {
		return nil, err
	}
	c.ContentSigning = *sec

	c.EntryCreditKey = factom.NewECAddress()

	return c, nil
}

func (a *AuthChannel) SignContent(msg []byte) []byte {
	return a.ContentSigning.Sign(msg)
}

func (a *AuthChannel) IsSameAs(b *AuthChannel) bool {
	if !a.Channel.IsSameAs(&b.Channel) {
		return false
	}

	for i := 0; i < PRIV_KEY_AMT; i++ {
		if !a.PrivateKeys[i].IsSameAs(&b.PrivateKeys[i]) {
			return false
		}
	}

	if !a.ContentSigning.IsSameAs(&b.ContentSigning) {
		return false
	}

	return true
}

func (a *AuthChannel) UnmarshalBinary(data []byte) error {
	_, err := a.UnmarshalBinaryData(data)
	return err
}

func (a *AuthChannel) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data

	newData, err = a.Channel.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	for i := 0; i < PRIV_KEY_AMT; i++ {
		newData, err = a.PrivateKeys[i].UnmarshalBinaryData(newData)
		if err != nil {
			return data, err
		}
	}

	newData, err = a.ContentSigning.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	if a.EntryCreditKey == nil {
		a.EntryCreditKey = new(factom.ECAddress)
	}

	newData, err = a.EntryCreditKey.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	return
}

func (a *AuthChannel) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data, err := a.Channel.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	for i := 0; i < PRIV_KEY_AMT; i++ {
		data, err = a.PrivateKeys[i].MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(data)
	}

	data, err = a.ContentSigning.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data = a.EntryCreditKey.SecBytes()[:32]
	buf.Write(data)

	return buf.Next(buf.Len()), nil
}
