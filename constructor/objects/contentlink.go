package objects

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	//"log"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
	"github.com/DistributedSolutions/DIMWIT/writeHelper/elements"
	log "github.com/DistributedSolutions/logrus"
)

// Factom Entry
//		byte		Version				0
//		[12]byte	"Content Link"		1
//		byte		ContentType			2
//		[32]byte	RootChainID			3
//		[32]byte	ContentChainID		4
//		[]byte		Timestamp			5
//		[32]byte	ContentSignKey		6
//		[64]byte	Signature			7
type ContentLinkApplyEntry struct {
	// Memory
	Channel *ChannelWrapper
	Entry   *lite.EntryHolder

	// Object
	Version       byte
	ContentType   byte
	TotalEntries  uint32
	RootChainID   primitives.Hash
	Infohash      primitives.InfoHash
	Timestamp     time.Time
	ShiftCipher   byte
	ContentSigKey primitives.PublicKey
	Message       []byte
	Signature     []byte

	// Content
	ContentChainID primitives.Hash
	ContentData    elements.ContentChainContent
	Content        common.Content
	LinkTimestamp  time.Time

	// Signal
	ErrorAndStop bool // Some verification failed, discard
}

func NewContentLinkApplyEntry() IApplyEntry {
	m := new(ContentLinkApplyEntry)
	m.ErrorAndStop = false
	return m
}

// Factom Entry
//		byte		Version				0
//		[12]byte	"Content Link"		1
//		byte		ContentType			2
//		[32]byte	RootChainID			3
// 		[32]byte 	ContentChain		4
//		[]byte		Timestamp			5
//		[32]byte	ContentSignKey		6
//		[64]byte	Signature			7
func (m *ContentLinkApplyEntry) ParseFactomEntry(e *lite.EntryHolder) error {
	m.Entry = e
	ex := e.Entry.ExtIDs

	r, err := primitives.BytesToHash(ex[3])
	if err != nil {
		return err
	}
	m.RootChainID = *r

	c, err := primitives.BytesToHash(ex[4])
	if err != nil {
		return err
	}
	m.ContentChainID = *c
	m.Content.ContentID = *c

	var t time.Time
	err = t.UnmarshalBinary(ex[5])
	if err != nil {
		return err
	}
	m.LinkTimestamp = t

	err = m.ContentSigKey.UnmarshalBinary(ex[6])
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	for i := 0; i < 6; i++ {
		buf.Write(ex[i])
	}
	m.Message = buf.Next(buf.Len())
	m.Signature = ex[7]

	// Do content when stiched

	return nil
}

func (m *ContentLinkApplyEntry) RequestChannel() (string, bool) {
	return m.RootChainID.String(), true
}

func (m *ContentLinkApplyEntry) AnswerChannelRequest(cw *ChannelWrapper) error {
	if cw == nil {
		return fmt.Errorf("Cannot work with nil channel")
	}
	m.Channel = cw
	return nil
}

func (m *ContentLinkApplyEntry) RequestEntriesInOtherChain() (string, bool) {
	return m.ContentChainID.String(), true
}

// Factom Chain
//		byte		Version				0
//		[13]byte	"Content Chain"		1
//		byte		ContentType			2
//		[4]byte		TotalEntries		3
//		[32]byte	RootChainID			4
//		[20]byte	Infohash			5
//		[]byte		Timestamp			6
//		byte		ShiftCipher			7
//		[32]byte	ContentSignKey		8
//		[64]byte	Signature			9
//		[]byte		nonce				10
//	CONTENT
//		XOR Marshaled content <-- Need to stitch
func (m *ContentLinkApplyEntry) AnswerChainEntriesInOther(first *lite.EntryHolder, rest []*lite.EntryHolder) {
	if bytes.Compare(first.Entry.ExtIDs[1], []byte("Content Chain")) != 0 {
		m.ErrorAndStop = true
	}

	ex := first.Entry.ExtIDs
	m.Version = ex[0][0]
	m.ContentType = ex[2][0]
	u, err := primitives.BytesToUint32(ex[3])
	if err != nil {
		log.Debug("[ContentLink] (1): Cannot unmarshal into uint32", err.Error())
		m.ErrorAndStop = true
		return
	}
	m.TotalEntries = u

	r, err := primitives.BytesToHash(ex[4])
	if err != nil {
		log.Debug("[ContentLink] (2): Cannot unmarshal into hash", err.Error())
		m.ErrorAndStop = true
		return
	}
	// Root chain in Link does not match.
	if !m.RootChainID.IsSameAs(r) {
		log.Debug("[ContentLink] (3): Root Chain ID does not match:")
		m.ErrorAndStop = true
		return
	}

	i, err := primitives.BytesToInfoHash(ex[5])
	if err != nil {
		log.Debug("[ContentLink] (4): Cannot unmarshal into infohash:", err.Error())
		m.ErrorAndStop = true
		return
	}
	m.Infohash = *i

	var t time.Time
	err = t.UnmarshalBinary(ex[6])
	if err != nil {
		log.Debug("[ContentLink] (5): Cannot unmarshal into time:", err.Error())
		m.ErrorAndStop = true
		return
	}
	m.Timestamp = t
	// Check Time window
	if !InsideTimeWindow(m.LinkTimestamp, m.Timestamp, constants.ENTRY_TIMESTAMP_WINDOW) {
		log.Debug("[ContentLink] (6): Time outside valid window")
		m.ErrorAndStop = true
		return
	}

	m.ShiftCipher = ex[7][0]

	err = m.ContentSigKey.UnmarshalBinary(ex[8])
	if err != nil {
		log.Debug("[ContentLink] (7): Cannot unmarshal into content signing key")
		m.ErrorAndStop = true
		return
	}

	buf := new(bytes.Buffer)
	for i := 0; i < 8; i++ {
		buf.Write(ex[i])
	}
	m.Message = buf.Next(buf.Len())
	m.Signature = ex[9]

	if !m.ContentSigKey.IsSameAs(&m.Channel.Channel.ContentSingingKey) {
		log.Debugf("[ContentLink] (8): Content signing key does not match. \n    Found %s, expect %s", m.ContentSigKey.String(), m.Channel.Channel.ContentSingingKey.String())
		m.ErrorAndStop = true
		return
	}

	if valid := m.ContentSigKey.Verify(m.Message, m.Signature); !valid {
		log.Debug("[ContentLink] (9): Content verify failed")
		m.ErrorAndStop = true
		return
	}

	// Stitch Content
	// Stich Entries :: 142bytes
	//		[4]byte		Sequence
	//		[32]byte	Sha256Hash of PreXOR
	//		[32]byte	ContentSignKey
	//		[64]byte	Signature
	// ContentStitch
	// data is the bytes to be stitched
	content := make([]byte, 0)
	plain := primitives.XORCipher(m.ShiftCipher, first.Entry.Content)
	content = append(content, plain...)

	var c uint32
	for c = 1; c < m.TotalEntries+1; c++ {
		// Looking for sequence c
		found := false
		for in, e := range rest {
			if len(e.Entry.ExtIDs) != 4 { // Crap
				//rest = append(rest[:in], rest[in+1:]...)
				continue
			}
			seq, err := primitives.BytesToUint32(e.Entry.ExtIDs[0])
			if err != nil {
				//rest = append(rest[:in], rest[in+1:]...) // Remove crap
				continue
			}
			if seq == uint32(c) {
				v, data := m.ValidateStitch(e)
				if v { // Stich applied, look for the next
					//rest = append(rest[:in], rest[in+1:]...)
					content = append(content, data...)
					found = true
					break
				}
			}
			var _ = in
		}
		if !found {
			log.Debug("[ContentLink] (10): Missing data in content link")
			m.ErrorAndStop = true
			return
		}
	}

	// Woo! Stiched up!
	err = m.ContentData.UnmarshalBinary(content)
	if err != nil {
		log.Debug("[ContentLink] (11): Cannot unmarshal into ContentData")
		m.ErrorAndStop = true
		return
	}

	m.Content = *m.ContentData.ContentChainContentToCommonContent()
	m.Content.Type = m.ContentType
	m.Content.CreationTime = m.LinkTimestamp
	m.Content.ContentID = m.ContentChainID
	m.RootChainID = m.RootChainID
}

func (m *ContentLinkApplyEntry) ValidateStitch(e *lite.EntryHolder) (bool, []byte) {
	plaintext := primitives.XORCipher(m.ShiftCipher, e.Entry.Content)
	sha := sha256.Sum256(plaintext)
	buf := new(bytes.Buffer)
	buf.Write(e.Entry.ExtIDs[0])
	buf.Write(sha[:])
	v := m.ContentSigKey.Verify(buf.Next(buf.Len()), e.Entry.ExtIDs[3])
	return v, plaintext
}

func (m *ContentLinkApplyEntry) ApplyEntry() (*ChannelWrapper, bool) {
	// Application happens above in the getting entries
	if m.ErrorAndStop {
		return nil, false
	}

	m.Content.RootChainID = m.RootChainID
	m.Channel.Channel.Content.ContentList = append(m.Channel.Channel.Content.ContentList, m.Content)

	return m.Channel, true
}

// unused
func (m *ContentLinkApplyEntry) NeedChainEntries() bool                      { return false }
func (m *ContentLinkApplyEntry) NeedIsFirstEntry() bool                      { return false }
func (m *ContentLinkApplyEntry) AnswerChainEntries(ents []*lite.EntryHolder) {}
func (m *ContentLinkApplyEntry) String() string                              { return "ContentLinkApplyEntry" }
