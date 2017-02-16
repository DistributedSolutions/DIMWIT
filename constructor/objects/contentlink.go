package objects

import (
	"time"

	"github.com/DistributedSolutions/DIMWIT/channelTool/creation"
	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

// Factom Chain
//		byte		Version				0
//		byte		ContentType			1
//		[4]byte		TotalEntries		2
//		[13]byte	"Content Chain"		3
//		[32]byte	RootChainID			4
//		[20]byte	Infohash			5
//		[]byte		Timestamp			6
//		byte		ShiftCipher			7
//		[32]byte	ContentSignKey		8
//		[64]byte	Signature			9
//		[]byte		nonce				10
//	CONTENT
//		Marshaled content <-- Need to stitch
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
	Signature     []byte

	// Content
	ContentData creation.ContentChainContent
	Content     common.Content
}

func NewContentLinkApplyEntry() IApplyEntry {
	m := new(ContentLinkApplyEntry)
	return m
}

func (m *ContentLinkApplyEntry) ParseFactomEntry(e *lite.EntryHolder) error {
	m.Entry = e
	ex := e.Entry.ExtIDs
	m.Version = ex[0][0]
	m.ContentType = ex[1][0]

	return nil
}

func (m *ContentLinkApplyEntry) RequestChannel() (string, bool)                     { return "", false }
func (m *ContentLinkApplyEntry) AnswerChannelRequest(cw *ChannelWrapper) error      { return nil }
func (m *ContentLinkApplyEntry) NeedChainEntries() bool                             { return false }
func (m *ContentLinkApplyEntry) NeedIsFirstEntry() bool                             { return false }
func (m *ContentLinkApplyEntry) AnswerChainEntries(ents []*lite.EntryHolder)        {}
func (m *ContentLinkApplyEntry) ApplyEntry() (*ChannelWrapper, bool)                { return nil, false }
func (m *ContentLinkApplyEntry) RequestEntriesInOtherChain() (string, bool)         { return "", false }
func (m *ContentLinkApplyEntry) AnswerChainEntriesInOther(ents []*lite.EntryHolder) {}
