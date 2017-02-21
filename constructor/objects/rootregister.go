package objects

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

// Factom entry
//		byte		Version
//		[13]byte	"Channel Chain"
//		[32]byte	RootChainID
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
type RootRegisterApplyEntry struct {
	// Memory
	Channel *ChannelWrapper
	Entry   *lite.EntryHolder

	//  Object
	RootChainID primitives.Hash
	Version     byte
	PubKey3     primitives.PublicKey
	Signature   []byte
	Message     []byte
}

func NewRootRegisterApplyEntry() IApplyEntry {
	m := new(RootRegisterApplyEntry)
	return m
}

func (m *RootRegisterApplyEntry) ParseFactomEntry(e *lite.EntryHolder) error {
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
	m.RootChainID = *hash
	return nil
}

func (m *RootRegisterApplyEntry) RequestChannel() (string, bool) {
	return m.RootChainID.String(), true
}

func (m *RootRegisterApplyEntry) AnswerChannelRequest(cw *ChannelWrapper) error {
	if cw == nil {
		return fmt.Errorf("Channel must exisit for RootRegisterApplyEntry")
	}
	m.Channel = cw
	return nil
}

func (m *RootRegisterApplyEntry) ApplyEntry() (*ChannelWrapper, bool) {
	if !m.PubKey3.IsSameAs(&m.Channel.Channel.LV3PublicKey) {
		return m.Channel, false // Invalid key
	}

	if valid := m.PubKey3.Verify(m.Message, m.Signature); !valid {
		return m.Channel, false // Bad signature
	}

	m.Channel.RRegistered = true
	return m.Channel, true
}

// Ununsed
func (m *RootRegisterApplyEntry) NeedChainEntries() bool                      { return false }
func (m *RootRegisterApplyEntry) NeedIsFirstEntry() bool                      { return false }
func (m *RootRegisterApplyEntry) AnswerChainEntries(ents []*lite.EntryHolder) {}
func (m *RootRegisterApplyEntry) RequestEntriesInOtherChain() (string, bool)  { return "", false }
func (m *RootRegisterApplyEntry) AnswerChainEntriesInOther(first *lite.EntryHolder, rest []*lite.EntryHolder) {
}
func (m *RootRegisterApplyEntry) String() string { return "RootRegisterApplyEntry" }
