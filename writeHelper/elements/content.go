package elements

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	//"encoding/hex"
	"fmt"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

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
//		Marshaled content

// Stich Entries :: 142bytes
//		[4]byte		Sequence
//		[32]byte	Sha256Hash of PreXOR
//		[32]byte	ContentSignKey
//		[64]byte	Signature
// ContentStitch
// data is the bytes to be stitched

type SingleContentChain struct {
	MetaData    ContentChainContent
	ContentKey  primitives.PrivateKey
	root        primitives.Hash
	cc          primitives.Hash
	contentType byte
}

func (SingleContentChain) Type() []byte  { return TYPE_SINGLE_CONTENT_CHAIN }
func (SingleContentChain) IsChain() bool { return true }
func (SingleContentChain) ForChain() int { return CHAIN_NA }

func (scc *SingleContentChain) Create(meta *ContentChainContent, contentKey primitives.PrivateKey, root primitives.Hash, contentChain primitives.Hash, contentType byte) {
	scc.MetaData = *meta
	scc.ContentKey = contentKey
	scc.cc = contentChain
	scc.root = root
	scc.contentType = contentType
}

func (scc *SingleContentChain) FactomElements() (*factom.Chain, []*factom.Entry, error) {
	placeHolder := make([]byte, 64)
	tsData, err := time.Now().MarshalBinary()
	if err != nil {
		return nil, nil, err
	}

	xorCipher := make([]byte, 1)
	rand.Read(xorCipher)
	xorKey := xorCipher[0]

	e := new(factom.Entry)
	extIDs := VersionAndType(scc)
	extIDs = append(extIDs, []byte{scc.contentType})       // 2
	extIDs = append(extIDs, primitives.Uint32ToBytes(0))   // 3
	extIDs = append(extIDs, scc.root.Bytes())              // 4
	extIDs = append(extIDs, scc.MetaData.InfoHash.Bytes()) // 5
	extIDs = append(extIDs, tsData)                        // 6
	extIDs = append(extIDs, []byte{xorKey})                // 7
	extIDs = append(extIDs, scc.ContentKey.Public.Bytes()) // 8
	extIDs = append(extIDs, placeHolder[:64])              // 9
	extIDs = append(extIDs, []byte{0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00}) // 10 Nonce holder

	// Find total entries needed
	headerLength := ExIDLength(extIDs)

	data, err := scc.MetaData.MarshalBinary()
	if err != nil {
		return nil, nil, err
	}
	contentLength := len(data)
	totalSize := contentLength + headerLength
	var entryCount int = 0
	if totalSize > constants.ENTRY_MAX_SIZE {
		entryCount = howManyEntries(headerLength, contentLength, 142)
	}
	extIDs[3] = primitives.Uint32ToBytes(uint32(entryCount))

	buf := new(bytes.Buffer)
	for i := 0; i < 8; i++ {
		buf.Write(extIDs[i])
	}
	msgData := buf.Next(buf.Len())
	sig := scc.ContentKey.Sign(msgData)
	extIDs[9] = sig

	nonce, chainID := FindValidNonce(extIDs[:10])
	extIDs[10] = nonce
	e.ExtIDs = extIDs

	firstEntContent := constants.ENTRY_MAX_SIZE - headerLength
	if entryCount > 0 {
		e.Content = primitives.XORCipher(xorKey, data[:firstEntContent])
		data = data[firstEntContent:]
	} else {
		e.Content = primitives.XORCipher(xorKey, data)
		data = []byte{}
	}

	c := factom.NewChain(e)
	// Entry 1 done

	// Stich Entries :: 142bytes
	//		[4]byte		Sequence
	//		[32]byte	Sha256Hash of PreXOR
	//		[32]byte	ContentSignKey
	//		[64]byte	Signature
	// ContentStitch
	// data is the bytes to be stitched
	bytesPerEntry := constants.ENTRY_MAX_SIZE - 142

	stiches := make([]*factom.Entry, entryCount)
	var seq uint32 = 1
	for len(data) > 0 {
		entry := new(factom.Entry)
		end := bytesPerEntry
		if len(data) < bytesPerEntry {
			end = len(data)
		}

		preXorContent := data[:end]
		entry.Content = primitives.XORCipher(xorKey, preXorContent)
		data = data[end:]

		hash := sha256.Sum256(preXorContent)
		// Set headers
		entry.ExtIDs = append(entry.ExtIDs, primitives.Uint32ToBytes(seq)) // 0
		entry.ExtIDs = append(entry.ExtIDs, hash[:])                       // 1

		msg := upToSig(entry.ExtIDs)
		entry.ExtIDs = append(entry.ExtIDs, scc.ContentKey.Public.Bytes()) // 2
		sig := scc.ContentKey.Sign(msg)
		entry.ExtIDs = append(entry.ExtIDs, sig) // 3
		hex, err := primitives.BytesToHash(chainID)
		if err != nil {
			return nil, nil, err
		}
		entry.ChainID = hex.String()

		if int(seq-1) >= len(stiches) {
			return nil, nil, fmt.Errorf("Ran out of entries. Seq is %d. Entrycount is %d, %d bytes left to write", seq-1, entryCount, len(data))
		}
		stiches[seq-1] = entry
		seq++
	}

	return c, stiches, nil
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
type ContentLinkEntry struct {
	rootChainID      primitives.Hash
	chanContentChain primitives.Hash
	contentChain     primitives.Hash
	contentKey       primitives.PrivateKey
	contentType      byte
}

func (ContentLinkEntry) Type() []byte  { return TYPE_CONTENT_LINK }
func (ContentLinkEntry) IsChain() bool { return false }
func (ContentLinkEntry) ForChain() int { return CHAIN_CONTENT_LIST }

func (cle *ContentLinkEntry) Create(root primitives.Hash, chanContentChain primitives.Hash, contentToLink primitives.Hash, contentKey primitives.PrivateKey, contentType byte) {
	cle.rootChainID = root
	cle.chanContentChain = chanContentChain
	cle.contentChain = contentToLink
	cle.contentKey = contentKey
	cle.contentType = contentType
}

func (cle *ContentLinkEntry) FactomEntry() (*factom.Entry, error) {
	tsData, err := time.Now().MarshalBinary()
	if err != nil {
		return nil, err
	}

	extIDs := VersionAndType(cle)
	extIDs = append(extIDs, []byte{cle.contentType})
	extIDs = append(extIDs, cle.rootChainID.Bytes())
	extIDs = append(extIDs, cle.contentChain.Bytes())
	extIDs = append(extIDs, tsData)

	sig := cle.contentKey.Sign(upToSig(extIDs))
	extIDs = append(extIDs, cle.contentKey.Public.Bytes())
	extIDs = append(extIDs, sig)

	e := new(factom.Entry)
	e.ExtIDs = extIDs
	e.ChainID = cle.chanContentChain.String()
	return e, nil
}

// The content of the entries in factom
type ContentChainContent struct {
	Title            primitives.Title
	LongDescription  primitives.LongDescription
	ShortDescription primitives.ShortDescription
	ActionFiles      primitives.FileList
	Thumbnail        primitives.Image
	Series           byte
	Part             [2]byte
	Tags             primitives.TagList
	// Torrent Metadata
	InfoHash     primitives.InfoHash
	Trackers     primitives.TrackerList
	TorrentFiles primitives.FileList
}

func RandomContentChainContent() *ContentChainContent {
	c := new(ContentChainContent)
	c.Title = *primitives.RandomTitle()
	c.LongDescription = *primitives.RandomLongDescription()
	c.ShortDescription = *primitives.RandomShortDescription()
	c.ActionFiles = *primitives.RandomFileList(uint32(10))
	c.Thumbnail = *primitives.RandomImage()
	c.Tags = *primitives.RandomTagList(uint32(constants.MAX_CONTENT_TAGS))
	c.InfoHash = *primitives.RandomInfoHash()
	c.Trackers = *primitives.RandomTrackerList(uint32(5))
	c.TorrentFiles = c.ActionFiles

	return c
}

func RandomHugeContentChainContent() *ContentChainContent {
	c := new(ContentChainContent)
	c.Title = *primitives.RandomTitle()
	c.LongDescription = *primitives.RandomLongDescription()
	c.ShortDescription = *primitives.RandomShortDescription()
	c.ActionFiles = *primitives.RandomFileList(uint32(10))
	c.Thumbnail = *primitives.RandomHugeImage()
	c.Tags = *primitives.RandomTagList(uint32(constants.MAX_CONTENT_TAGS))
	c.InfoHash = *primitives.RandomInfoHash()
	c.Trackers = *primitives.RandomTrackerList(uint32(5))
	c.TorrentFiles = c.ActionFiles

	return c
}

func (c *ContentChainContent) ContentChainContentToCommonContent() *common.Content {
	co := new(common.Content)
	co.ContentTitle = c.Title
	co.LongDescription = c.LongDescription
	co.ShortDescription = c.ShortDescription
	co.ActionFiles = c.ActionFiles
	co.Thumbnail = c.Thumbnail
	co.Series = c.Series
	co.Part = c.Part
	co.Tags = c.Tags
	co.InfoHash = c.InfoHash
	co.Trackers = c.Trackers
	co.FileList = c.TorrentFiles
	return co
}

func CommonContentToContentChainContent(c *common.Content) *ContentChainContent {
	co := new(ContentChainContent)
	co.Title = c.ContentTitle
	co.LongDescription = c.LongDescription
	co.ShortDescription = c.ShortDescription
	co.ActionFiles = c.ActionFiles
	co.Thumbnail = c.Thumbnail
	co.Series = c.Series
	co.Part = c.Part
	co.Tags = c.Tags
	co.InfoHash = c.InfoHash
	co.Trackers = c.Trackers
	co.TorrentFiles = c.FileList
	return co
}

func (c *ContentChainContent) MarshalBinary() (data []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[Content] A panic has occurred while marshaling: %s", r)
			return
		}
	}()

	buf := new(bytes.Buffer)

	data, err = c.Title.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.LongDescription.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.ShortDescription.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.ActionFiles.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.Thumbnail.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	buf.Write([]byte{c.Series})

	buf.Write(c.Part[:])

	data, err = c.Tags.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	// Torrent Metadata

	data, err = c.InfoHash.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.Trackers.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.TorrentFiles.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	return buf.Next(buf.Len()), err
}

func (p *ContentChainContent) UnmarshalBinary(data []byte) error {
	_, err := p.UnmarshalBinaryData(data)
	return err
}

func (c *ContentChainContent) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data

	newData, err = c.Title.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.LongDescription.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.ShortDescription.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.ActionFiles.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.Thumbnail.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	c.Series = newData[0]
	newData = newData[1:]

	copy(c.Part[:], newData[:2])
	newData = newData[2:]

	newData, err = c.Tags.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.InfoHash.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.Trackers.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.TorrentFiles.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	return
}

func (a *ContentChainContent) IsSameAs(b *ContentChainContent) bool {
	if !a.Title.IsSameAs(&b.Title) {
		return false
	}

	if !a.LongDescription.IsSameAs(&b.LongDescription) {
		return false
	}

	if !a.ShortDescription.IsSameAs(&b.ShortDescription) {
		return false
	}

	if !a.ActionFiles.IsSameAs(&b.ActionFiles) {
		return false
	}

	if !a.Thumbnail.IsSameAs(&b.Thumbnail) {
		return false
	}

	if a.Series != b.Series {
		return false
	}

	if a.Part[0] != b.Part[0] || a.Part[1] != b.Part[1] {
		return false
	}

	if !a.Tags.IsSameAs(&b.Tags) {
		return false
	}

	if !a.InfoHash.IsSameAs(&b.InfoHash) {
		return false
	}

	if !a.Trackers.IsSameAs(&b.Trackers) {
		return false
	}

	if !a.TorrentFiles.IsSameAs(&b.TorrentFiles) {
		return false
	}

	return true
}
