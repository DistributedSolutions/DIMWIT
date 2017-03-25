// The tool to create and save channels and their private keys. Also covers content creation
package channelTool

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

type CreationTool struct {
	Channels map[string]*AuthChannel
}

func NewCreationTool() *CreationTool {
	ct := new(CreationTool)
	ct.Channels = make(map[string]*AuthChannel)

	return ct
}

func (ct *CreationTool) AddNewChannel(ch *common.Channel,
	contentSigKey primitives.PrivateKey,
	filePath string,
	ec *factom.ECAddress) error {
	if _, ok := ct.Channels[ch.RootChainID.String()]; ok {
		return fmt.Errorf("Channel already exists in the CreationTool")
	}

	a, err := NewAuthChannel(ch, ec)
	if err != nil {
		return err
	}

	a.TorrentUploadPath.SetString(filePath)

	ct.Channels[a.Channel.RootChainID.String()] = a

	return nil
}

// AddPrivateKey adds a private key to authority channel
// 3 == content signing key
func (ct *CreationTool) AddPrivateKey(lvl int, key primitives.PrivateKey, root primitives.Hash) error {
	if lvl > 2 || lvl < 0 {
		return fmt.Errorf("lvl is out of range. Expect 0-3, found %d", lvl)
	}

	a, ok := ct.Channels[root.String()]
	if !ok || a == nil {
		return fmt.Errorf("channel not found")
	}

	if lvl == 3 {
		ct.Channels[root.String()].ContentSigning = key
	} else {
		ct.Channels[root.String()].PrivateKeys[lvl] = key
	}
	return nil
}

func (ct *CreationTool) SetUploadPath(path string, root primitives.Hash) error {
	a, ok := ct.Channels[root.String()]
	if !ok || a == nil {
		return fmt.Errorf("channel not found")
	}

	return ct.Channels[root.String()].TorrentUploadPath.SetString(path)
}

// CreateAllFactomEntries should be called when you first make a channel
func (ct *CreationTool) CreateAllFactomEntries(root primitives.Hash) error {
	a, ok := ct.Channels[root.String()]
	if !ok || a == nil {
		return fmt.Errorf("channel not found")
	}

	if a.Channel.Status() < constants.CHANNEL_READY {
		return fmt.Errorf("Channel given is not ready, it is missing elements")
	}

	if !factom.IsValidAddress(a.EntryCreditKey.String()) {
		return fmt.Errorf("Entry credit address is invalid")
	}

	err := a.MakeChannel()
	if err != nil {
		return err
	}

	err = a.MakeManagerChain()
	if err != nil {
		return err
	}

	err = a.MakeContentChain()
	if err != nil {
		return err
	}

	err = a.MakeContents()
	if err != nil {
		return err
	}

	return nil
}

func (ct *CreationTool) AddContent(root primitives.Hash,
	con *common.Content) (chains []*factom.Chain,
	ents []*factom.Entry,
	err error) {

	a, ok := ct.Channels[root.String()]
	if !ok || a == nil {
		return nil, nil, fmt.Errorf("channel not found")
	}

	cc, err := a.AddContent(con)
	if err != nil {
		return nil, nil, err
	}

	return cc.ReturnChains(), cc.ReturnEntries(), nil
}
