package objects

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/DistributedSolutions/DIMWIT/channelTool/creation"
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
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
	Meta     creation.ManageChainMetaData
	DBlockTS time.Time
}

func NewManageMetaApplyEntry() IApplyEntry {
	m := new(ManageMetaApplyEntry)
	meta := new(creation.ManageChainMetaData)
	meta.ChannelTags = primitives.NewTagList(uint32(constants.MAX_CHANNEL_TAGS))
	m.Meta = *meta
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
	m.Timestamp = ts

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
	if !m.RootChain.IsSameAs(&m.Channel.Channel.RootChainID) {
		return nil, false // Wrong chain dumbass
	}

	if !m.PubKey3.IsSameAs(&m.Channel.Channel.LV3PublicKey) {
		return nil, false // Invalid key
	}

	if valid := m.PubKey3.Verify(m.Message, m.Signature); !valid {
		return nil, false // Bad sig
	}

	hash := sha256.Sum256(m.Entry.Entry.Content)
	if bytes.Compare(hash[:], m.Entry.Entry.ExtIDs[5]) != 0 {
		return nil, false // Bad content hash
	}

	if m.Entry.Entry.ChainID != m.Channel.Channel.ManagementChainID.String() {
		return nil, false // Entry in wrong chain
	}

	m.DBlockTS = time.Unix(m.Entry.Timestamp, 0)
	if !InsideTimeWindow(m.DBlockTS, m.Timestamp, constants.ENTRY_TIMESTAMP_WINDOW) {
		return nil, false // Too old
	}

	// Need to stich all the content
	content := make([]byte, 0)
	content = append(content, m.Entry.Entry.Content...)

	rest := m.Entries
	var c uint32
	for c = 1; c < m.EntryCount+1; c++ {
		// Looking for sequence c
		found := false
		for in, e := range rest {
			var _ = in
			if len(e.Entry.ExtIDs) != 9 ||
				bytes.Compare(e.Entry.ExtIDs[1], []byte("Channel Management Metadata Stich")) != 0 { // Crap
				//rest = RemoveFromList(rest, in) // Remove crap
				//skip = true
				//fmt.Println("REMOVE", e.Entry.ExtIDs[4], len(rest))
				continue
			}

			u, err := primitives.BytesToUint32(e.Entry.ExtIDs[4])
			if err != nil {
				//rest = RemoveFromList(rest, in) // Remove crap
				//fmt.Println("REMOVE", e.Entry.ExtIDs[4], len(rest))
				continue
			}

			seq := u
			if seq == c {
				v, data := m.ValidateStitch(e)
				if v { // Stich applied, look for the next
					//rest = RemoveFromList(rest, in)
					content = append(content, data...)
					found = true
					break
				}
			}
		}

		if !found {
			return nil, false
		}
	}

	fContHash := sha256.Sum256(content)
	if bytes.Compare(fContHash[:], m.FullContentHash.Bytes()) != 0 {
		return nil, false
	}

	t := new(creation.ManageChainMetaData)
	m.Meta = *t

	// Wow, we did it.
	err := m.Meta.UnmarshalBinary(content)
	if err != nil {
		return nil, false // Fuuuuuck
	}

	ch := m.Channel.Channel
	ch = metaToChannel(ch, m.Meta)
	m.Channel.Channel = ch

	return m.Channel, true
}

func metaToChannel(ch common.Channel, meta creation.ManageChainMetaData) common.Channel {
	if meta.Website != nil {
		ch.Website = *meta.Website
	}
	if meta.LongDescription != nil {
		ch.LongDescription = *meta.LongDescription
	}
	if meta.ShortDescription != nil {
		ch.ShortDescription = *meta.ShortDescription
	}
	if meta.Playlist != nil {
		ch.Playlist = *(ch.Playlist.Combine(meta.Playlist))
		//ch.Playlist = meta.Playlist
	}
	if meta.Thumbnail != nil {
		ch.Thumbnail = *meta.Thumbnail
	}
	if meta.Banner != nil {
		ch.Banner = *meta.Banner
	}
	if meta.ChannelTags != nil {
		ch.Tags = *ch.Tags.Combine(meta.ChannelTags)
	}
	if meta.SuggestedChannels != nil {
		ch.SuggestedChannel = *ch.SuggestedChannel.Combine(meta.SuggestedChannels)
	}
	return ch
}

// Entry Stich
//	0	byte		Version
//	1	[33]byte	"Channel Management Metadata Stich"
//	2	[32]byte	RootChainID
//	3	[32]byte	FullContentHash
//	4	[4]byte		Sequence
//	5	[32]byte	ContentHash
//	6	[15]byte	Timestamp
//	7	[32]byte	PublicKey(3)
//	8	[64]byte	Signature
func (m *ManageMetaApplyEntry) ValidateStitch(e *lite.EntryHolder) (bool, []byte) {
	hash, err := primitives.BytesToHash(e.Entry.ExtIDs[2])
	if err != nil {
		return false, nil
	}
	if !hash.IsSameAs(&m.RootChain) {
		return false, nil
	}

	hash, err = primitives.BytesToHash(e.Entry.ExtIDs[3])
	if err != nil {
		return false, nil
	}
	if !hash.IsSameAs(&m.FullContentHash) {
		return false, nil
	}

	hash, err = primitives.BytesToHash(e.Entry.ExtIDs[5])
	if err != nil {
		return false, nil
	}

	contentHash := sha256.Sum256(e.Entry.Content)
	if bytes.Compare(contentHash[:], hash.Bytes()) != 0 {
		return false, nil
	}

	var ts time.Time
	err = ts.UnmarshalBinary(e.Entry.ExtIDs[6])
	if err != nil {
		return false, nil
	}
	if !InsideTimeWindow(m.DBlockTS, ts, constants.ENTRY_TIMESTAMP_WINDOW) {
		return false, nil // Too old
	}

	if bytes.Compare(m.PubKey3.Bytes(), e.Entry.ExtIDs[7]) != 0 {
		return false, nil // Wrong key
	}

	buf := new(bytes.Buffer)
	for i := 0; i < 7; i++ {
		buf.Write(e.Entry.ExtIDs[i])
	}
	msg := buf.Next(buf.Len())
	if valid := m.PubKey3.Verify(msg, e.Entry.ExtIDs[8]); !valid {
		return false, nil
	}

	return true, e.Entry.Content
}

// unused
func (m *ManageMetaApplyEntry) RequestEntriesInOtherChain() (string, bool) { return "", false }
func (m *ManageMetaApplyEntry) AnswerChainEntriesInOther(first *lite.EntryHolder, rest []*lite.EntryHolder) {
}
func (m *ManageMetaApplyEntry) String() string { return "ManageMetaApplyEntry" }
