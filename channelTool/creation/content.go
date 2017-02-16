package creation

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

type ContentChain struct {
	FirstEntry *factom.Chain
	Entries    []*factom.Entry

	// Must be done after first entry
	Register *RegisterStruct
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
//		Marshaled content

//  CreateContentChain needs all the metadata to determine how many entries to use
func (r *ContentChain) CreateContentChain(contentType byte, contentData ContentChainContent, root primitives.Hash, contentSignKey primitives.PrivateKey) (err error) {
	/*defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("A panic has occurred while creating content chain: %s", r)
			return
		}
	}()*/

	// HEADER
	timeData, err := time.Now().MarshalBinary()
	if err != nil {
		return err
	}

	xorCipher := make([]byte, 1)
	rand.Read(xorCipher)
	xorKey := xorCipher[0]

	e := new(factom.Entry)
	e.ExtIDs = append(e.ExtIDs, []byte{constants.FACTOM_VERSION}) // 0
	e.ExtIDs = append(e.ExtIDs, []byte{contentType})              // 1
	e.ExtIDs = append(e.ExtIDs, primitives.Uint32ToBytes(0))      // 2
	e.ExtIDs = append(e.ExtIDs, []byte("Content Chain"))          // 3
	e.ExtIDs = append(e.ExtIDs, root.Bytes())                     // 4
	e.ExtIDs = append(e.ExtIDs, contentData.InfoHash.Bytes())     // 5
	e.ExtIDs = append(e.ExtIDs, timeData)                         // 6
	e.ExtIDs = append(e.ExtIDs, []byte{xorKey})                   // 7
	e.ExtIDs = append(e.ExtIDs, contentSignKey.Public.Bytes())    // 8

	msg := upToNonce(e.ExtIDs, 8)
	sig := contentSignKey.Sign(msg)
	e.ExtIDs = append(e.ExtIDs, sig) // 9
	e.ExtIDs = append(e.ExtIDs, []byte{0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00}) // 10 Nonce holder

	// Find total entries needed
	headerLength := ExIDLength(e.ExtIDs)

	data, err := contentData.MarshalBinary()
	if err != nil {
		return err
	}
	contentLength := len(data)
	totalSize := contentLength + headerLength
	var entryCount int = 0
	if totalSize > constants.ENTRY_MAX_SIZE {
		entryCount = howManyEntries(headerLength, contentLength, 142)
	}
	e.ExtIDs[2] = primitives.Uint32ToBytes(uint32(entryCount))

	// Find nonce
	c := new(CreateStruct)
	c.endExtID = 10
	c.ExtIDs = e.ExtIDs[:10]
	nonce := FindValidNonce(c)
	e.ExtIDs[10] = nonce // 10

	firstEntContent := constants.ENTRY_MAX_SIZE - headerLength

	if entryCount > 0 {
		e.Content = primitives.XORCipher(xorKey, data[:firstEntContent])
		data = data[firstEntContent:]
	} else {
		e.Content = primitives.XORCipher(xorKey, data)
		data = []byte{}
	}

	r.FirstEntry = factom.NewChain(e)

	// Stich Entries :: 142bytes
	//		[4]byte		Sequence
	//		[32]byte	Sha256Hash of PreXOR
	//		[32]byte	ContentSignKey
	//		[64]byte	Signature
	// ContentStitch
	// data is the bytes to be stitched

	r.Entries = make([]*factom.Entry, entryCount)
	bytesPerEntry := constants.ENTRY_MAX_SIZE - 142
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
		entry.ExtIDs = append(entry.ExtIDs, contentSignKey.Public.Bytes()) // 2

		msg := upToNonce(entry.ExtIDs, 3)
		sig := contentSignKey.Sign(msg)
		entry.ExtIDs = append(entry.ExtIDs, sig) // 4
		entry.ChainID = r.FirstEntry.ChainID

		if int(seq-1) >= len(r.Entries) {
			return fmt.Errorf("Ran out of entries. Seq is %d. Entrycount is %d, %d bytes left to write", seq-1, entryCount, len(data))
		}
		r.Entries[seq-1] = entry
		seq++
	}

	return nil
}

// Factom Entry
//		byte		Version
//		byte		ContentType
//		[12]byte	"Content Link"
//		[32]byte	RootChainID
//		[]byte		Timestamp
//		[32]byte	ContentSignKey
//		[64]byte	Signature
func (r *ContentChain) RegisterNewContentChain(rootChain primitives.Hash, contentChainID primitives.Hash, contentType byte, sigKey primitives.PrivateKey) error {
	e := new(factom.Entry)

	timeData, err := time.Now().MarshalBinary()
	if err != nil {
		return fmt.Errorf("Unable to create a timestamp: %s", err.Error())
	}

	e.ExtIDs = append(e.ExtIDs, []byte{constants.FACTOM_VERSION}) // 0
	e.ExtIDs = append(e.ExtIDs, []byte{contentType})              // 1
	e.ExtIDs = append(e.ExtIDs, []byte("Content Link"))           // 2
	e.ExtIDs = append(e.ExtIDs, rootChain.Bytes())                // 3
	e.ExtIDs = append(e.ExtIDs, timeData)                         // 4
	e.ExtIDs = append(e.ExtIDs, sigKey.Public.Bytes())            // 5

	msg := upToNonce(e.ExtIDs, 4)
	sig := sigKey.Sign(msg)
	e.ExtIDs = append(e.ExtIDs, sig) // 4

	e.ChainID = contentChainID.String()
	r.Register.Entry = e

	return nil
}

func howManyEntries(headerLength int, contentLength int, contentHeaderLength int) int {
	contentLength -= (constants.ENTRY_MAX_SIZE - headerLength)
	bytesPerEntry := constants.ENTRY_MAX_SIZE - contentHeaderLength
	count := 0
	for contentLength > 0 {
		contentLength -= bytesPerEntry
		count++
	}

	return count
}
