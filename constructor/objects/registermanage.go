package objects

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

// Factom Entry
//		byte		Version
//		[25]byte	"Register Management Chain"
//		[32]byte	ManagementChainID
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
type ManageRegisterApplyEntry struct {
	// Memory
	Channel *ChannelWrapper
	Entry   *lite.EntryHolder

	//  Object
	ManageChain primitives.Hash
	Version     byte
	PubKey3     primitives.PublicKey
	Signature   []byte
	Message     []byte
}

func NewManageRegisterApplyEntry() IApplyEntry {
	m := new(ManageRegisterApplyEntry)
	return m
}

func (m *ManageRegisterApplyEntry) ParseFactomEntry(e *lite.EntryHolder) error {
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
	m.ManageChain = *hash
	return nil
}

func (m *ManageRegisterApplyEntry) RequestChannel() (string, bool) {
	return m.Entry.Entry.ChainID, true
}

func (m *ManageRegisterApplyEntry) AnswerChannelRequest(cw *ChannelWrapper) error {
	if cw == nil {
		return fmt.Errorf("Channel must exisit for RootRegisterApplyEntry")
	}
	m.Channel = cw
	return nil
}

func (m *ManageRegisterApplyEntry) ApplyEntry() (*ChannelWrapper, bool) {
	if !m.PubKey3.IsSameAs(&m.Channel.Channel.LV3PublicKey) {
		return m.Channel, false // Invalid key
	}

	if valid := m.PubKey3.Verify(m.Message, m.Signature); !valid {
		return m.Channel, false // Bad signature
	}

	if !m.Channel.Channel.ManagementChainID.IsSameAs(&m.ManageChain) {
		return m.Channel, false
	}

	if m.Entry.Entry.ChainID != m.Channel.Channel.RootChainID.String() {
		return nil, false // Must be in root
	}

	m.Channel.MRegistered = true
	return m.Channel, true
}

// Ununsed
func (m *ManageRegisterApplyEntry) NeedChainEntries() bool                      { return false }
func (m *ManageRegisterApplyEntry) NeedIsFirstEntry() bool                      { return false }
func (m *ManageRegisterApplyEntry) AnswerChainEntries(ents []*lite.EntryHolder) {}
func (m *ManageRegisterApplyEntry) RequestEntriesInOtherChain() (string, bool)  { return "", false }
func (m *ManageRegisterApplyEntry) AnswerChainEntriesInOther(first *lite.EntryHolder, rest []*lite.EntryHolder) {
}
func (m *ManageRegisterApplyEntry) String() string { return "ManageRegisterApplyEntry" }
