package objects

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

// Factom Entry
//		byte		Version
//		[25]byte	"Register Content Chain"
//		[32]byte	Channel Content ChainID
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
type ContentRegisterApplyEntry struct {
	// Memory
	Channel *ChannelWrapper
	Entry   *lite.EntryHolder

	//  Object
	ContentChain primitives.Hash
	Version      byte
	PubKey3      primitives.PublicKey
	Signature    []byte
	Message      []byte
}

func NewContentRegisterApplyEntry() IApplyEntry {
	m := new(ContentRegisterApplyEntry)
	return m
}

func (m *ContentRegisterApplyEntry) ParseFactomEntry(e *lite.EntryHolder) error {
	m.Entry = e

	m.Version = e.Entry.ExtIDs[0][0]

	buf := new(bytes.Buffer)
	buf.Write(e.Entry.ExtIDs[0])
	buf.Write(e.Entry.ExtIDs[1])
	buf.Write(e.Entry.ExtIDs[2])
	m.Message = buf.Next(buf.Len())

	err := m.PubKey3.UnmarshalBinary(e.Entry.ExtIDs[3])
	if err != nil {
		return err
	}

	m.Signature = e.Entry.ExtIDs[4]

	hash, err := primitives.BytesToHash(e.Entry.ExtIDs[2])
	if err != nil {
		return err
	}
	m.ContentChain = *hash
	return nil
}

func (m *ContentRegisterApplyEntry) RequestChannel() (string, bool) {
	return m.Entry.Entry.ChainID, true
}

func (m *ContentRegisterApplyEntry) AnswerChannelRequest(cw *ChannelWrapper) error {
	if cw == nil {
		return fmt.Errorf("Channel must exisit for RootRegisterApplyEntry")
	}
	m.Channel = cw
	return nil
}

func (m *ContentRegisterApplyEntry) ApplyEntry() (*ChannelWrapper, bool) {
	if !m.PubKey3.IsSameAs(&m.Channel.Channel.LV3PublicKey) {
		return nil, false // Invalid key
	}

	if valid := m.PubKey3.Verify(m.Message, m.Signature); !valid {
		return nil, false // Bad signature
	}

	if !m.Channel.Channel.ContentChainID.IsSameAs(&m.ContentChain) {
		return nil, false
	}

	if m.Entry.Entry.ChainID != m.Channel.Channel.RootChainID.String() {
		return nil, false // Must be in root
	}

	m.Channel.CRegistered = true
	return m.Channel, true
}

// Ununsed
func (m *ContentRegisterApplyEntry) NeedChainEntries() bool                      { return false }
func (m *ContentRegisterApplyEntry) NeedIsFirstEntry() bool                      { return false }
func (m *ContentRegisterApplyEntry) AnswerChainEntries(ents []*lite.EntryHolder) {}
func (m *ContentRegisterApplyEntry) RequestEntriesInOtherChain() (string, bool)  { return "", false }
func (m *ContentRegisterApplyEntry) AnswerChainEntriesInOther(first *lite.EntryHolder, rest []*lite.EntryHolder) {
}
