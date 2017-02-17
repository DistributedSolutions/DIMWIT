package objects

import (
	"bytes"
	"fmt"
	"time"

	"github.com/DistributedSolutions/DIMWIT/channelTool/creation"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

type ManageMetaApplyEntry struct {
	// Memory
	Channel *ChannelWrapper
	Entry   *lite.EntryHolder
	Entries []*lite.EntryHolder

	// Object
	Version         byte
	RootChain       primitives.Hash
	FullContentHash primitives.Hash
	EntryCount      uint32
	Timestamp       time.Time
	PubKey3         primitives.PublicKey
	Message         []byte
	Signature       []byte

	// Holders
	Meta creation.ManageChainMetaData
}

func NewManageMetaApplyEntry() IApplyEntry {
	m := new(ManageMetaApplyEntry)
	return m
}

// All entries simply overwrite
// Entry Main
//	0	byte		Version
//	1	[32]byte	"Channel Management Metadata Main"
//	2	[32]byte	RootChainID
//	3	[32]byte	FullContentHash
//	4	[32]byte	EntryCount
//	5	[32]byte	ContentHash
//	6	[15]byte	Timestamp
//	7	[32]byte	PublicKey(3)
//	8	[64]byte	Signature
func (m *ManageMetaApplyEntry) ParseFactomEntry(e *lite.EntryHolder) error {
	ex := e.Entry.ExtIDs
	m.Version = ex[0][0]

	hash, err := primitives.BytesToHash(ex[2])
	if err != nil {
		return err
	}
	m.RootChain = *hash

	hash, err = primitives.BytesToHash(ex[3])
	if err != nil {
		return err
	}
	m.FullContentHash = *hash

	u, err := primitives.BytesToUint32(ex[4])
	if err != nil {
		return err
	}
	m.EntryCount = u

	var ts time.Time
	err = ts.UnmarshalBinary(ex[6])
	if err != nil {
		return err
	}

	err = m.PubKey3.UnmarshalBinary(ex[7])
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	for i := 0; i < 7; i++ {
		buf.Write(e.Entry.ExtIDs[i])
	}
	m.Message = buf.Next(buf.Len())
	m.Signature = e.Entry.ExtIDs[8]
	m.Entry = e
	return nil
}

func (m *ManageMetaApplyEntry) RequestChannel() (string, bool) {
	return m.RootChain.String(), true
}

func (m *ManageMetaApplyEntry) AnswerChannelRequest(cw *ChannelWrapper) error {
	if cw == nil {
		return fmt.Errorf("Channel must exisit for RootRegisterApplyEntry")
	}
	m.Channel = cw
	return nil
}

func (m *ManageMetaApplyEntry) NeedChainEntries() bool {
	return true
}

func (m *ManageMetaApplyEntry) NeedIsFirstEntry() bool {
	return false
}

func (m *ManageMetaApplyEntry) AnswerChainEntries(ents []*lite.EntryHolder) {
	m.Entries = ents
}

func (m *ManageMetaApplyEntry) ApplyEntry() (*ChannelWrapper, bool) {

	return nil, false
}

// unused
func (m *ManageMetaApplyEntry) RequestEntriesInOtherChain() (string, bool) { return "", false }
func (m *ManageMetaApplyEntry) AnswerChainEntriesInOther(first *lite.EntryHolder, rest []*lite.EntryHolder) {
}

// Data from Entries in MsgPack
//		Channel Website
//		Channel LongDescription
//		Channel ShortDescription
//		Playlist
//		Thumbnail
//		Banner
//		Channel Tags
//		Suggested Channels

// Entry Stich
//	0	byte		Version
//	1	[33]byte	"Channel Management Metadata Stich"
//	2	[32]byte	RootChainID
//	3	[32]byte	FullContentHash
//	4	[4]byte		Sequence
//	5	[32]byte	ContentHash
//	6	[15]byte		Timestamp
//	7	[32]byte	PublicKey(3)
//	8	[64]byte	Signature
