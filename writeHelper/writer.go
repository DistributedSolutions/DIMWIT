package writeHelper

import (
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/constructor"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
	"github.com/DistributedSolutions/DIMWIT/util"
	"github.com/DistributedSolutions/DIMWIT/writeHelper/elements"
	"github.com/FactomProject/factom"
)

type WriteHelper struct {
	// To write into Factom
	Writer lite.FactomLiteWriter

	// To read from Factom
	Reader *constructor.Constructor

	// Map of AuthChannels
	AuthChannels map[string]*AuthChannel

	// ECKey
	ECAddress *factom.ECAddress
}

func NewWriterHelper(con *constructor.Constructor, fw lite.FactomLiteWriter) (*WriteHelper, error) {
	w := new(WriteHelper)
	w.Reader = con
	w.Writer = fw

	w.AuthChannels = make(map[string]*AuthChannel)

	pk, err := primitives.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	w.ECAddress, err = factom.MakeECAddress(pk.Secret[:32])
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *WriteHelper) SetECAddress(sec string) error {
	ec, err := factom.GetECAddress(sec)
	if err != nil {
		return err
	}
	w.ECAddress = ec
	return nil
}

func (w *WriteHelper) GetECAddress() *factom.ECAddress {
	return w.ECAddress
}

func (w *WriteHelper) MakeNewAuthChannel(ch *common.Channel) error {
	_, err := w.InitiateChannel(ch)
	if err != nil {
		return err
	}
	err = w.UpdateChannel(ch)
	if err != nil {
		return err
	}

	for i, c := range ch.Content.ContentList {
		c.RootChainID = ch.RootChainID
		p := &c
		_, err := w.AddContent(p)
		if err != nil {
			return err
		}
		ch.Content.ContentList[i] = *p
	}

	return nil
}

type CostStruct struct {
	InitCost   int `json:"initcost"`
	UpdateCost int `json:"updatecost"`
}

func (w *WriteHelper) VerifyChannel(ch *common.Channel) (cost *CostStruct, apiErr *util.ApiError) {
	cs := new(CostStruct)
	// Init cost
	// 	3 Chains = 30
	// 	3 Register = 3
	cs.InitCost = 33
	cs.UpdateCost = -1

	metaData := w.createMetaDataChanges(ch, nil)
	data, err := metaData.MarshalBinary()
	if err != nil {
		return cs, util.NewAPIErrorFromOne(err)
	}
	// Cost of the Management MetaData
	ec := elements.HowManyEntries(elements.ManageContentHeaderLength, len(data), 248)
	cs.UpdateCost = ec

	return cs, nil
}

func (w *WriteHelper) InitiateChannel(ch *common.Channel) (*primitives.Hash, *util.ApiError) {
	// TODO: Check Balance

	// Generate Keys
	a, err := NewAuthChannel(ch, w.ECAddress)
	if err != nil {
		return nil, util.NewAPIError(err, fmt.Errorf("failed to generate channel keys"))
	}

	// Brute force a ChainIDs
	a.Initiate(ch)

	// Get Factom Elements
	//		Root
	root := a.Root
	entries, chain := root.FactomElements()

	//	Enter into Factom
	w.Writer.SubmitChain(*chain, *w.ECAddress)
	for _, e := range entries {
		w.Writer.SubmitEntry(*e, *w.ECAddress)
	}

	//		Manage
	manage := a.Manage
	entries, chain = manage.FactomElements()

	//	Enter into Factom
	w.Writer.SubmitChain(*chain, *w.ECAddress)
	for _, e := range entries {
		w.Writer.SubmitEntry(*e, *w.ECAddress)
	}

	//		ContentList
	cc := a.ContentList
	entries, chain = cc.FactomElements()

	//	Enter into Factom
	w.Writer.SubmitChain(*chain, *w.ECAddress)
	for _, e := range entries {
		w.Writer.SubmitEntry(*e, *w.ECAddress)
	}

	// Add to our Map
	w.AuthChannels[ch.RootChainID.String()] = a
	return &a.ChannelRoot, nil
}

func (w *WriteHelper) UpdateChannel(ch *common.Channel) (apiErr *util.ApiError) {
	a, ok := w.AuthChannels[ch.RootChainID.String()]
	if !ok || a == nil {
		util.NewAPIErrorFromOne(fmt.Errorf("We do not have the keys to this channel"))
	}

	factomChannel, err := w.Reader.RetrieveChannel(a.ChannelRoot)
	if err != nil {
		return util.NewAPIErrorFromOne(err)
	}

	var metaData *elements.ManageChainMetaData
	if factomChannel == nil {
		metaData = w.createMetaDataChanges(ch, nil)
	} else {
		metaData = w.createMetaDataChanges(ch, &(factomChannel.Channel))
	}

	a.Manage.MetaData.Create(metaData, a.PrivateKeys[2], a.ChannelRoot, a.ChannelManage)
	ents, err := a.Manage.MetaData.FactomEntry()
	if err != nil {
		return util.NewAPIErrorFromOne(err)
	}

	for _, e := range ents {
		w.Writer.SubmitEntry(*e, *w.ECAddress)
	}

	return
}

func (w *WriteHelper) AddExistingChannel(pk *primitives.PublicKey) (err *util.ApiError) {
	return &util.ApiError{
		LogError:  fmt.Errorf("Cannot add existing channel"),
		UserError: fmt.Errorf("Cannot add existing channel"),
	}
}

func (w *WriteHelper) DeleteChannel(rootChain *primitives.Hash) (apiErr *util.ApiError) {
	return &util.ApiError{
		LogError:  fmt.Errorf("Cannot delete channel"),
		UserError: fmt.Errorf("Cannot delete channel"),
	}
}

func (w *WriteHelper) VerifyContent(ch *common.Content) (cost int, apiErr *util.ApiError) {
	return 0, util.NewAPIError(nil, nil)
}

func (w *WriteHelper) AddContent(con *common.Content) (hash *primitives.Hash, apiErr *util.ApiError) {
	a, ok := w.AuthChannels[con.RootChainID.String()]
	if !ok {
		return nil, util.NewAPIErrorFromOne(fmt.Errorf("Do not have the keys for that channel"))
	}

	con.Thumbnail = *primitives.RandomHugeImage()
	ce := new(elements.SingleContentChain)
	ce.Create(elements.CommonContentToContentChainContent(con), a.ContentSigning, a.ChannelRoot, a.ChannelContent, con.Type)
	c, entries, err := ce.FactomElements()
	if err != nil {
		return nil, util.NewAPIErrorFromOne(err)
	}

	w.Writer.SubmitChain(*c, *w.ECAddress)
	he, _ := primitives.HexToHash(c.FirstEntry.ChainID)
	con.ContentID = *he

	for _, e := range entries {
		w.Writer.SubmitEntry(*e, *w.ECAddress)
	}
	return &con.ContentID, nil
}

func (w *WriteHelper) DeleteContent(contentID *primitives.Hash) (apiErr *util.ApiError) {
	return &util.ApiError{
		LogError:  fmt.Errorf("Cannot delete Content"),
		UserError: fmt.Errorf("Cannot delete Content"),
	}
}

/*
type ManageChainMetaData struct {
	Website           *primitives.SiteURL
	LongDescription   *primitives.LongDescription
	ShortDescription  *primitives.ShortDescription
	Playlist          *common.ManyPlayList
	Thumbnail         *primitives.Image
	Banner            *primitives.Image
	ChannelTags       *primitives.TagList
	SuggestedChannels *primitives.HashList
}
*/

func newele() *elements.ManageChainMetaData {
	ele := new(elements.ManageChainMetaData)
	ele.Website = new(primitives.SiteURL)
	ele.LongDescription = new(primitives.LongDescription)
	ele.ShortDescription = new(primitives.ShortDescription)
	ele.Playlist = new(common.ManyPlayList)
	ele.Thumbnail = new(primitives.Image)
	ele.Banner = new(primitives.Image)
	ele.ChannelTags = new(primitives.TagList)
	ele.SuggestedChannels = new(primitives.HashList)

	return ele
}

func (w *WriteHelper) createMetaDataChanges(ch *common.Channel, factomChannel *common.Channel) *elements.ManageChainMetaData {
	changes := newele()
	if factomChannel == nil { // Only New
		if !ch.Website.Empty() && changes.Website != nil {
			*changes.Website = ch.Website
		}

		if changes.LongDescription != nil && !ch.LongDescription.Empty() {
			*changes.LongDescription = ch.LongDescription
		}

		if changes.ShortDescription != nil && !ch.ShortDescription.Empty() {
			*changes.ShortDescription = ch.ShortDescription
		}

		if changes.Playlist != nil && !ch.Playlist.Empty() {
			*changes.Playlist = ch.Playlist
		}

		if changes.Thumbnail != nil && !ch.Thumbnail.Empty() {
			*changes.Thumbnail = ch.Thumbnail
		}

		if changes.Banner != nil && !ch.Banner.Empty() {
			*changes.Banner = ch.Banner
		}

		if changes.ChannelTags != nil && !ch.Tags.Empty() {
			*changes.ChannelTags = ch.Tags
		}

		if changes.SuggestedChannels != nil && !ch.SuggestedChannel.Empty() {
			*changes.SuggestedChannels = ch.SuggestedChannel
		}
	} else { // Compare
		if !ch.Website.IsSameAs(&factomChannel.Website) {
			*changes.Website = ch.Website
		}

		if !ch.LongDescription.IsSameAs(&factomChannel.LongDescription) {
			*changes.LongDescription = ch.LongDescription
		}

		if !ch.ShortDescription.IsSameAs(&factomChannel.ShortDescription) {
			*changes.ShortDescription = ch.ShortDescription
		}

		if !ch.Playlist.IsSameAs(&factomChannel.Playlist) {
			*changes.Playlist = ch.Playlist
		}

		if !ch.Thumbnail.IsSameAs(&factomChannel.Thumbnail) {
			*changes.Thumbnail = ch.Thumbnail
		}

		if !ch.Banner.IsSameAs(&factomChannel.Banner) {
			*changes.Banner = ch.Banner
		}

		if !ch.Tags.IsSameAs(&factomChannel.Tags) {
			*changes.ChannelTags = ch.Tags
		}

		if !ch.SuggestedChannel.IsSameAs(&factomChannel.SuggestedChannel) {
			*changes.SuggestedChannels = ch.SuggestedChannel
		}
	}

	return changes
}
