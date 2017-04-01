// The tool to create and save channels and their private keys. Also covers content creation
package channelTool

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factom/wallet"
)

type CreationTool struct {
	Channels map[string]*AuthChannel
	Wallet   *wallet.Wallet
}

func NewCreationTool() (*CreationTool, error) {
	ct := new(CreationTool)
	var err error

	ct.Channels = make(map[string]*AuthChannel)
	// TOOD: Make a saved wallet, right now live in a map
	ct.Wallet, err = wallet.NewMapDBWallet()
	if err != nil {
		return nil, err
	}

	return ct, nil
}

func (ct *CreationTool) GetECAddress(root primitives.Hash) (*factom.ECAddress, error) {
	a, ok := ct.Channels[root.String()]
	if !ok || a == nil {
		return nil, fmt.Errorf("channel not found")
	}

	return a.EntryCreditKey, nil
}

func (ct *CreationTool) ReturnFactomElements(root primitives.Hash) ([]*factom.Entry, []*factom.Chain, error) {
	a, ok := ct.Channels[root.String()]
	if !ok || a == nil {
		return nil, nil, fmt.Errorf("channel not found")
	}

	if a.RootChain == nil {
		err := a.MakeChannel()
		if err != nil {
			return nil, nil, err
		}
	}

	if a.ManageChain == nil {
		err := a.MakeManagerChain()
		if err != nil {
			return nil, nil, err
		}
	}

	if a.ContentChain == nil {
		err := a.MakeContentChain()
		if err != nil {
			return nil, nil, err
		}
	}

	chains, err := a.ReturnFactomChains()
	if err != nil {
		return nil, nil, err
	}

	ents, err := a.ReturnFactomEntries()
	if err != nil {
		return nil, nil, err
	}
	return ents, chains, nil
}

func (ct *CreationTool) AddNewChannel(ch *common.Channel, filePaths []string) (*primitives.Hash, error) {
	if _, ok := ct.Channels[ch.RootChainID.String()]; ok {
		return nil, fmt.Errorf("Channel already exists in the CreationTool")
	}

	// TODO: Refactor to have master EC address
	tempEc, err := ct.Wallet.GenerateECAddress()
	if err != nil {
		return nil, err
	}

	a, err := MakeNewAuthChannel(ch, tempEc)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(filePaths); i++ {
		a.TorrentUploadPaths[i].SetString(filePaths[i])
	}

	ct.Channels[a.Channel.RootChainID.String()] = a

	return &a.Channel.RootChainID, nil
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

func (ct *CreationTool) SetUploadPath(path string, root primitives.Hash, contentIndex uint) error {
	a, ok := ct.Channels[root.String()]
	if !ok || a == nil {
		return fmt.Errorf("channel not found")
	}

	return ct.Channels[root.String()].TorrentUploadPaths[contentIndex].SetString(path)
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
	con *common.Content) (ents []*factom.Entry,
	chains []*factom.Chain,
	err error) {

	a, ok := ct.Channels[root.String()]
	if !ok || a == nil {
		return nil, nil, fmt.Errorf("channel not found")
	}

	cc, err := a.AddContent(con)
	if err != nil {
		return nil, nil, err
	}

	return cc.ReturnEntries(), cc.ReturnChains(), nil
}
