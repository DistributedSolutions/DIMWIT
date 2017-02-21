package objects

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

// Factom Chain
//		byte		Version
//		[24]byte	"Channel Content Chain"
//		[32]byte	RootChainID
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
//		[]byte		nonce
type ContentChainApplyEntry struct {
	// Memory
	Channel *ChannelWrapper
	Entry   *lite.EntryHolder

	//  Object
	Version     byte
	RootChainID primitives.Hash
	Signature   []byte
	Message     []byte
	PubKey3     primitives.PublicKey
}

func NewContentChainApplyEntry() IApplyEntry {
	m := new(ContentChainApplyEntry)
	return m
}

func (r *ContentChainApplyEntry) ParseFactomEntry(e *lite.EntryHolder) error {
	ent := e.Entry
	r.Version = ent.ExtIDs[0][0]
	err := r.PubKey3.UnmarshalBinary(ent.ExtIDs[3])
	if err != nil {
		return err
	}

	err = r.RootChainID.UnmarshalBinary(ent.ExtIDs[2])
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.Write(e.Entry.ExtIDs[0])
	buf.Write(e.Entry.ExtIDs[1])
	buf.Write(e.Entry.ExtIDs[2])
	r.Message = buf.Next(buf.Len())
	r.Signature = ent.ExtIDs[4]

	r.Entry = e
	return nil
}

func (r *ContentChainApplyEntry) RequestChannel() (string, bool) {
	return r.RootChainID.String(), true
}

func (r *ContentChainApplyEntry) AnswerChannelRequest(cw *ChannelWrapper) error {
	if cw == nil {
		return fmt.Errorf("Channel must exisit for ContentChainApplyEntry")
	}
	r.Channel = cw
	return nil
}

func (r *ContentChainApplyEntry) NeedIsFirstEntry() bool { return true }

func (m *ContentChainApplyEntry) ApplyEntry() (*ChannelWrapper, bool) {
	if !m.PubKey3.IsSameAs(&m.Channel.Channel.LV3PublicKey) {
		return m.Channel, false // Invalid key
	}

	if valid := m.PubKey3.Verify(m.Message, m.Signature); !valid {
		return m.Channel, false // Bad signature
	}

	hash, err := primitives.HexToHash(m.Entry.Entry.ChainID)
	if err != nil {
		return nil, false
	}

	m.Channel.CMadeHeight = m.Entry.Height
	m.Channel.Channel.ContentChainID = *hash
	return m.Channel, true
}

// Unused
func (r *ContentChainApplyEntry) NeedChainEntries() bool                      { return false }
func (r *ContentChainApplyEntry) AnswerChainEntries(ents []*lite.EntryHolder) {}
func (m *ContentChainApplyEntry) RequestEntriesInOtherChain() (string, bool)  { return "", false }
func (m *ContentChainApplyEntry) AnswerChainEntriesInOther(first *lite.EntryHolder, rest []*lite.EntryHolder) {
}
func (m *ContentChainApplyEntry) String() string { return "ContentChainApplyEntry" }
