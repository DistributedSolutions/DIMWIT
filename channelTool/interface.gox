package channelTool

import (
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/util"
	"github.com/FactomProject/factom"
)

type IChannelTool interface {
	// Verifies this channel can be added
	VerifyChannel(ch *common.Channel) (cost int, err *util.ApiError)

	// Generate private keys and store. Also sets RootchainID, Manage, and ContentID in channel
	// Only set title and keys. Will not set any metadata. TOTALLY EMPTY!! (Aside from title)
	InitiateChannel(ch *common.Channel) (chains []*factom.Chain, entries []*factom.Entry, err *util.ApiError)

	// Create factom elements to apply changes to channel in blockchain. Will create managment
	// chain if needed. Does not do content
	UpdateChannel(ch *common.Channel) (chains []*factom.Chain, entries []*factom.Entry, err *util.ApiError)

	// Will delete channel given root chain id. Will check to see if we have the correct keys.
	DeleteChannel(rootChain *primitives.Hash) (err *util.ApiError)

	// Verifies this content can be added
	VerifyContent(ch *common.Content) (cost int, err *util.ApiError)

	// Will add content given root chain id. Will check to see if we have the correct keys.
	AddContent(con *common.Content, contentID *primitives.Hash) (chains []*factom.Chain, entries []*factom.Entry, err *util.ApiError)

	// Will delete content given root chain id. Will check to see if we have the correct keys.
	DeleteContent(contentID *primitives.Hash) (err *util.ApiError)
}

// A fake thing to simulate a tool
type FakeChannelTool struct {
}

func NewFakeChannelTool() IChannelTool {
	return new(FakeChannelTool)
}

func (fakeChannelTool *FakeChannelTool) VerifyChannel(ch *common.Channel) (cost int, err *util.ApiError) {
	return 0, util.NewAPIError(nil, nil)
}

func (fakeChannelTool *FakeChannelTool) InitiateChannel(ch *common.Channel) (chains []*factom.Chain, entries []*factom.Entry, err *util.ApiError) {
	return
}

func (fakeChannelTool *FakeChannelTool) UpdateChannel(ch *common.Channel) (chains []*factom.Chain, entries []*factom.Entry, err *util.ApiError) {
	return
}

func (fakeChannelTool *FakeChannelTool) DeleteChannel(rootChain *primitives.Hash) (err *util.ApiError) {
	return
}

func (fakeChannelTool *FakeChannelTool) VerifyContent(ch *common.Content) (cost int, err *util.ApiError) {
	return 0, util.NewAPIError(nil, nil)
}

func (fakeChannelTool *FakeChannelTool) AddContent(con *common.Content, contentID *primitives.Hash) (chains []*factom.Chain, entries []*factom.Entry, err *util.ApiError) {
	return
}

func (fakeChannelTool *FakeChannelTool) DeleteContent(contentID *primitives.Hash) (err *util.ApiError) {
	return
}
