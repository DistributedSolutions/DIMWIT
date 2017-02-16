package objects

import (
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

// ApplyEntry is capable of taking a factom entry,
// and apply the factom entry to a channel.
type IApplyEntry interface {
	ParseFactomEntry(e *lite.EntryHolder) error
	RequestChannel() (string, bool)
	AnswerChannelRequest(cw *ChannelWrapper)

	// Special
	NeedChainEntries() bool
	NeedIsFirstEntry() bool
	AnswerChainEntries(ents []*lite.EntryHolder)
	AnswerFirstEntry(a bool)

	// ApplyEntry returns the channel and a bool to indicate wether or not
	// it made changes
	ApplyEntry() (*ChannelWrapper, bool)
}
