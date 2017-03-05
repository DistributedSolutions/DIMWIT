package objects

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
	log "github.com/DistributedSolutions/logrus"
)

//// Factom Chain
//		byte		Version
//		[24]byte	"Channel Management Chain"
//		[32]byte	RootChainID
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
//		[]byte		nonce
type ManageChainApplyEntry struct {
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

func NewManageChainApplyEntry() IApplyEntry {
	m := new(ManageChainApplyEntry)
	return m
}

func (r *ManageChainApplyEntry) ParseFactomEntry(e *lite.EntryHolder) error {
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

func (r *ManageChainApplyEntry) RequestChannel() (string, bool) {
	return r.RootChainID.String(), true
}

func (r *ManageChainApplyEntry) AnswerChannelRequest(cw *ChannelWrapper) error {
	if cw == nil {
		return fmt.Errorf("Channel must exist for ManageChainApplyEntry")
	}
	r.Channel = cw
	return nil
}

func (r *ManageChainApplyEntry) NeedIsFirstEntry() bool { return true }

func (m *ManageChainApplyEntry) ApplyEntry() (*ChannelWrapper, bool) {
	if !m.PubKey3.IsSameAs(&m.Channel.Channel.LV3PublicKey) {
		log.Debug("[ManageChain] (1): Public key does not match")
		return m.Channel, false // Invalid key
	}

	if valid := m.PubKey3.Verify(m.Message, m.Signature); !valid {
		log.Debug("[ManageChain] (2): Bad signature")
		return m.Channel, false // Bad signature
	}

	hash, err := primitives.HexToHash(m.Entry.Entry.ChainID)
	if err != nil {
		log.Debug("[ManageChain] (3): Cannot unmarshal into hex")
		return nil, false
	}

	m.Channel.MMadeHeight = m.Entry.Height
	m.Channel.Channel.ManagementChainID = *hash
	return m.Channel, true
}

// Unused
func (r *ManageChainApplyEntry) NeedChainEntries() bool                      { return false }
func (r *ManageChainApplyEntry) AnswerChainEntries(ents []*lite.EntryHolder) {}
func (m *ManageChainApplyEntry) RequestEntriesInOtherChain() (string, bool)  { return "", false }
func (m *ManageChainApplyEntry) AnswerChainEntriesInOther(first *lite.EntryHolder, rest []*lite.EntryHolder) {
}
func (m *ManageChainApplyEntry) String() string { return "ManageChainApplyEntry" }
