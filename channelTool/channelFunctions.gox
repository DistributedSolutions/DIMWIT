package channelTool

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/channelTool/creation"
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

// Order:
// 	MakeChannel
//  MakeManagerChain
//  MakeContentChain

func (a *AuthChannel) ReturnFactomChains() ([]*factom.Chain, error) {
	if a.RootChain == nil || a.ManageChain == nil || a.ContentChain == nil {
		return nil, fmt.Errorf("Missing chains with true: \nRoot - %t\nManage - %t\nContent - %t\n",
			a.RootChain == nil, a.ManageChain == nil, a.ContentChain == nil)
	}

	c := make([]*factom.Chain, 0)
	c = append(c, a.RootChain.ReturnChains()...)
	c = append(c, a.ManageChain.ReturnChains()...)
	c = append(c, a.ContentChain.ReturnChains()...)

	for _, con := range a.Contents {
		c = append(c, con.ReturnChains()...)
	}
	return c, nil
}

func (a *AuthChannel) ReturnFactomEntries() ([]*factom.Entry, error) {
	if a.RootChain == nil || a.ManageChain == nil || a.ContentChain == nil {
		return nil, fmt.Errorf("Missing chains with true: \nRoot - %t\nManage - %t\nContent - %t\n",
			a.RootChain == nil, a.ManageChain == nil, a.ContentChain == nil)
	}

	c := make([]*factom.Entry, 0)
	c = append(c, a.RootChain.ReturnEntries()...)
	c = append(c, a.ManageChain.ReturnEntries()...)
	c = append(c, a.ContentChain.ReturnEntries()...)

	for _, con := range a.Contents {
		c = append(c, con.ReturnEntries()...)
	}

	return c, nil
}

// Makes the root chain and registers into the master
func (a *AuthChannel) MakeChannel() error {
	rc := new(creation.RootChain)

	pubs := make([]primitives.PublicKey, 3)
	for i, p := range a.PrivateKeys {
		pubs[i] = p.Public
	}

	err := rc.CreateRootChain(pubs, a.Channel.ChannelTitle)
	if err != nil {
		return err
	}

	h, err := primitives.HexToHash(rc.Create.Chain.ChainID)
	if err != nil {
		return err
	}
	a.Channel.RootChainID = *h

	rc.RegisterRootEntry(a.Channel.RootChainID, a.PrivateKeys[2])

	err = rc.ContentSigningKey(a.Channel.RootChainID, a.Channel.ContentSingingKey, a.PrivateKeys[2])
	if err != nil {
		return err
	}

	a.RootChain = rc
	return nil
}

func (a *AuthChannel) MakeManagerChain() error {
	mc := new(creation.ManageChain)
	if a.Channel.RootChainID.IsSameAs(primitives.NewZeroHash()) {
		return fmt.Errorf("No root chain found")
	}

	err := mc.CreateManagementChain(a.Channel.RootChainID, a.PrivateKeys[2])
	if err != nil {
		return err
	}

	h, err := primitives.HexToHash(mc.Create.Chain.FirstEntry.ChainID)
	if err != nil {
		return err
	}
	a.Channel.ManagementChainID = *h

	meta := creation.NewManageChainMetaData()
	meta.Website = &a.Channel.Website
	meta.LongDescription = &a.Channel.LongDescription
	meta.ShortDescription = &a.Channel.ShortDescription
	meta.Playlist = &a.Channel.Playlist
	meta.Thumbnail = &a.Channel.Thumbnail
	meta.Banner = &a.Channel.Banner
	meta.ChannelTags = &a.Channel.Tags
	meta.SuggestedChannels = &a.Channel.SuggestedChannel

	err = mc.CreateMetadata(meta, a.Channel.RootChainID, a.Channel.ManagementChainID, a.PrivateKeys[2])
	if err != nil {
		return err
	}

	mc.RegisterChannelManagementChain(a.Channel.RootChainID, a.Channel.ManagementChainID, a.PrivateKeys[2])
	a.ManageChain = mc
	return nil
}

func (a *AuthChannel) MakeContentChain() error {
	cc := new(creation.ChanContentChain)
	err := cc.CreateChanContentChain(a.Channel.RootChainID, a.PrivateKeys[2])
	if err != nil {
		return err
	}

	h, err := primitives.HexToHash(cc.Create.Chain.FirstEntry.ChainID)
	if err != nil {
		return err
	}
	a.Channel.ContentChainID = *h

	cc.RegisterChannelContentChain(a.Channel.RootChainID, *h, a.PrivateKeys[2])

	a.ContentChain = cc
	return nil
}

func (a *AuthChannel) MakeContents() error {
	for i, c := range a.Channel.Content.GetContents() {
		cc := new(creation.ContentChain)
		cont := creation.CommonContentToContentChainContent(&c)
		err := cc.CreateContentChain(c.Type, *cont, a.Channel.RootChainID, a.ContentSigning)
		if err != nil {
			return err
		}

		cc.RegisterNewContentChain(a.Channel.RootChainID, a.Channel.ContentChainID, c.Type, a.ContentSigning)
		a.Contents = append(a.Contents, cc)
		a.Channel.Content.ContentList[i].RootChainID = a.Channel.RootChainID
		chainID, err := primitives.HexToHash(cc.FirstEntry.FirstEntry.ChainID)
		if err != nil {
			return err
		}
		a.Channel.Content.ContentList[i].ContentID = *chainID
	}

	return nil
}

func (a *AuthChannel) AddContent(c *common.Content) (*creation.ContentChain, error) {
	cc := new(creation.ContentChain)
	cont := creation.CommonContentToContentChainContent(c)
	err := cc.CreateContentChain(c.Type, *cont, a.Channel.RootChainID, a.ContentSigning)
	if err != nil {
		return nil, err
	}

	cc.RegisterNewContentChain(a.Channel.RootChainID, a.Channel.ContentChainID, c.Type, a.ContentSigning)
	a.Contents = append(a.Contents, cc)

	c.RootChainID = a.Channel.RootChainID
	chainID, err := primitives.HexToHash(cc.FirstEntry.FirstEntry.ChainID)
	if err != nil {
		return nil, err
	}

	c.ContentID = *chainID
	a.Channel.Content.ContentList = append(a.Channel.Content.ContentList, *c)
	return cc, nil
}
