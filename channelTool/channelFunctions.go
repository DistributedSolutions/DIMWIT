package channelTool

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/channelTool/creation"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

// Order:
// 	Register channel
//

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

	h, err := primitives.HexToHash(rc.Create.Chain.FirstEntry.ChainID)
	if err != nil {
		return err
	}
	a.Channel.RootChainID = *h

	rc.RegisterRootEntry(a.Channel.RootChainID, a.PrivateKeys[2])

	a.RootChain = rc
	return nil
}

func (a *AuthChannel) ReturnFactomChains() []*factom.Chain {
	c := make([]*factom.Chain, 0)
	c = append(c, a.RootChain.ReturnChains()...)
	c = append(c, a.ManageChain.ReturnChains()...)
	c = append(c, a.ContentChain.ReturnChains()...)
	return c
}

func (a *AuthChannel) ReturnFactomEntries() []*factom.Entry {
	c := make([]*factom.Entry, 0)
	c = append(c, a.RootChain.ReturnEntries()...)
	c = append(c, a.ManageChain.ReturnEntries()...)
	c = append(c, a.ContentChain.ReturnEntries()...)
	return c
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
