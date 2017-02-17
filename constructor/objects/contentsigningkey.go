package objects

import (
	"bytes"
	"fmt"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

type ContentSigningKeyApplyEntry struct {
	// Memory
	Channel *ChannelWrapper
	Entry   *lite.EntryHolder

	//  Object
	Version       byte
	RootChain     primitives.Hash
	ContentSigKey primitives.PublicKey
	TimeStamp     time.Time
	PubKey3       primitives.PublicKey
	Signature     []byte
	Message       []byte
}

func NewContentSigningKeyApplyEntry() IApplyEntry {
	m := new(ContentSigningKeyApplyEntry)
	return m
}

// Factom entry
//		byte		Version
//		[19]byte	"Content Signing Key"
//		[32]byte	RootChainID
//		[32]byte	ContentSigningKey
//		[15]byte	Timestamp
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
func (m *ContentSigningKeyApplyEntry) ParseFactomEntry(e *lite.EntryHolder) error {
	m.Entry = e

	m.Version = e.Entry.ExtIDs[0][0]

	buf := new(bytes.Buffer)
	for i := 0; i < 5; i++ {
		buf.Write(e.Entry.ExtIDs[i])
	}
	m.Message = buf.Next(buf.Len())

	hash, err := primitives.BytesToHash(e.Entry.ExtIDs[2])
	if err != nil {
		return err
	}
	m.RootChain = *hash

	err = m.ContentSigKey.UnmarshalBinary(e.Entry.ExtIDs[3])
	if err != nil {
		return err
	}

	var ts time.Time
	err = ts.UnmarshalBinary(e.Entry.ExtIDs[4])
	if err != nil {
		return err
	}
	m.TimeStamp = ts

	err = m.PubKey3.UnmarshalBinary(e.Entry.ExtIDs[5])
	if err != nil {
		return err
	}

	m.Signature = e.Entry.ExtIDs[6]
	return nil
}

func (m *ContentSigningKeyApplyEntry) RequestChannel() (string, bool) {
	return m.Entry.Entry.ChainID, true
}

func (m *ContentSigningKeyApplyEntry) AnswerChannelRequest(cw *ChannelWrapper) error {
	if cw == nil {
		return fmt.Errorf("Channel must exisit for RootRegisterApplyEntry")
	}
	m.Channel = cw
	return nil
}

func (m *ContentSigningKeyApplyEntry) ApplyEntry() (*ChannelWrapper, bool) {
	if !m.PubKey3.IsSameAs(&m.Channel.Channel.LV3PublicKey) {
		return m.Channel, false // Invalid key
	}

	if valid := m.PubKey3.Verify(m.Message, m.Signature); !valid {
		return m.Channel, false // Bad signature
	}

	if !m.Channel.Channel.RootChainID.IsSameAs(&m.RootChain) {
		return m.Channel, false
	}

	if m.RootChain.String() != m.Entry.Entry.ChainID {
		return m.Channel, false // In the wrong chain
	}

	etime := time.Unix(m.Entry.Timestamp, 0)
	if !InsideTimeWindow(etime, m.TimeStamp, constants.ENTRY_TIMESTAMP_WINDOW) {
		return m.Channel, false // Bad timestamp
	}

	m.Channel.Channel.ContentSingingKey = m.ContentSigKey
	return m.Channel, true
}

// Ununsed
func (m *ContentSigningKeyApplyEntry) NeedChainEntries() bool                      { return false }
func (m *ContentSigningKeyApplyEntry) NeedIsFirstEntry() bool                      { return false }
func (m *ContentSigningKeyApplyEntry) AnswerChainEntries(ents []*lite.EntryHolder) {}
func (m *ContentSigningKeyApplyEntry) RequestEntriesInOtherChain() (string, bool)  { return "", false }
func (m *ContentSigningKeyApplyEntry) AnswerChainEntriesInOther(first *lite.EntryHolder, rest []*lite.EntryHolder) {
}
