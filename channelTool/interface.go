package channelTool

import (
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

type IChannelTool interface {
	// Verifies this channel can be added
	VerifyChannel(ch *common.Channel) (cost int, err error)

	// Generate private keys and store. Also sets RootchainID, Manage, and ContentID in channel
	// Only set title and keys. Will not set any metadata. TOTALLY EMPTY!! (Aside from title)
	InitiateChannel(ch *common.Channel) (chains []*factom.Chain, entries []*factom.Entry, err error)

	// Create factom elements to apply changes to channel in blockchain. Will create managment
	// chain if needed. Does not do content
	UpdateChannel(ch *common.Channel) (chains []*factom.Chain, entries []*factom.Entry, err error)

	// Will add content given root chain id. Will check to see if we have the correct keys.
	AddContent(con *common.Content, rootChain *primitives.Hash) (chains []*factom.Chain, entries []*factom.Entry, err error)
}

// A fake thing to simulate a tool
type FakeChannelTool struct {
}

func NewFakeChannelTool() IChannelTool {
	return new(FakeChannelTool)
}

func VerifyChannel(ch *common.Channel) (cost int, err error) {
	return
}

func InitiateChannel(ch *common.Channel) (chains []*factom.Chain, entries []*factom.Entry, err error) {
	return
}

func UpdateChannel(ch *common.Channel) (chains []*factom.Chain, entries []*factom.Entry, err error) {
	return
}

func AddContent(con *common.Content, rootChain *primitives.Hash) (chains []*factom.Chain, entries []*factom.Entry, err error) {
	return
}
