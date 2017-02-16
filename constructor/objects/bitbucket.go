package objects

import (
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

// BitBucketApplyEntry doesn't process or do anything
// It's a way to process an entry, but do nothing
type BitBucketApplyEntry struct{}

func NewBitBucketApplyEntry() IApplyEntry {
	m := new(BitBucketApplyEntry)
	return m
}
func (m *BitBucketApplyEntry) ParseFactomEntry(e *lite.EntryHolder) error         { return nil }
func (m *BitBucketApplyEntry) RequestChannel() (string, bool)                     { return "", false }
func (m *BitBucketApplyEntry) AnswerChannelRequest(cw *ChannelWrapper) error      { return nil }
func (m *BitBucketApplyEntry) NeedChainEntries() bool                             { return false }
func (m *BitBucketApplyEntry) NeedIsFirstEntry() bool                             { return false }
func (m *BitBucketApplyEntry) AnswerChainEntries(ents []*lite.EntryHolder)        {}
func (m *BitBucketApplyEntry) ApplyEntry() (*ChannelWrapper, bool)                { return nil, false }
func (m *BitBucketApplyEntry) RequestEntriesInOtherChain() (string, bool)         { return "", false }
func (m *BitBucketApplyEntry) AnswerChainEntriesInOther(ents []*lite.EntryHolder) {}
