package objects

import (
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

// This one was easy.
type MasterChainApplyEntry struct{}

func NewMasterChainApplyEntry() IApplyEntry {
	m := new(MasterChainApplyEntry)
	return m
}
func (m *MasterChainApplyEntry) ParseFactomEntry(e *lite.EntryHolder) error    { return nil }
func (m *MasterChainApplyEntry) RequestChannel() (string, bool)                { return "", false }
func (m *MasterChainApplyEntry) AnswerChannelRequest(cw *ChannelWrapper) error { return nil }
func (m *MasterChainApplyEntry) NeedChainEntries() bool                        { return false }
func (m *MasterChainApplyEntry) NeedIsFirstEntry() bool                        { return false }
func (m *MasterChainApplyEntry) AnswerChainEntries(ents []*lite.EntryHolder)   {}
func (m *MasterChainApplyEntry) ApplyEntry() (*ChannelWrapper, bool)           { return nil, false }
func (m *MasterChainApplyEntry) RequestEntriesInOtherChain() (string, bool)    { return "", false }
func (m *MasterChainApplyEntry) AnswerChainEntriesInOther(first *lite.EntryHolder, rest []*lite.EntryHolder) {
}
