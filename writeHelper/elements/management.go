package elements

import (
	"crypto/sha256"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/FactomProject/factom"
)

type Manage struct {
	ManageChain         *ManageChain
	RegisterManageChain *RegisterManageChain
	MetaData            *ManageMetaData
}

func NewManage() *Manage {
	m := new(Manage)
	m.ManageChain = new(ManageChain)
	m.RegisterManageChain = new(RegisterManageChain)
	m.MetaData = new(ManageMetaData)

	return m
}

// Factom Chain
//		byte		Version
//		[24]byte	"Channel Management Chain"
//		[32]byte	RootChainID
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
//		[]byte		nonce
type ManageChain struct {
	RootChainID primitives.Hash
	KeyToSign   primitives.PrivateKey
	Nonce       []byte
}

func (m *ManageChain) Create(root primitives.Hash, key3 primitives.PrivateKey) (*primitives.Hash, error) {
	m.RootChainID = root

	nonce, chainID := FindValidNonce(m.AllButNonce())
	m.Nonce = nonce

	return primitives.BytesToHash(chainID)
}

func (m *ManageChain) AllButNonce() [][]byte {
	extIDs := VersionAndType(m)
	extIDs = append(extIDs, m.RootChainID.Bytes())

	sig := m.KeyToSign.Sign(upToSig(extIDs))
	extIDs = append(extIDs, m.KeyToSign.Public.Bytes())
	extIDs = append(extIDs, sig)

	return extIDs
}

func (m *ManageChain) FactomChain() *factom.Chain {
	e := new(factom.Entry)
	extIDs := m.AllButNonce()
	e.ExtIDs = append(extIDs, m.Nonce)

	return factom.NewChain(e)
}

func (ManageChain) Type() []byte  { return TYPE_MANAGE_CHAIN }
func (ManageChain) IsChain() bool { return true }
func (ManageChain) ForChain() int { return CHAIN_NA }

// Factom Entry
//		byte		Version
//		[25]byte	"Register Management Chain"
//		[32]byte	ManagementChainID
//		[32]byte	PublicKey(3)
//		[64]byte	Signature
type RegisterManageChain struct {
	RootChainID   primitives.Hash
	ManageChainID primitives.Hash
	KeyToSign     primitives.PrivateKey
}

func (RegisterManageChain) Type() []byte  { return TYPE_MANAGE_CHAIN_REGISTER }
func (RegisterManageChain) IsChain() bool { return false }
func (RegisterManageChain) ForChain() int { return CHAIN_ROOT }

func (rmc *RegisterManageChain) Create(rootChain primitives.Hash, manageChainId primitives.Hash, key3 primitives.PrivateKey) {
	rmc.RootChainID = rootChain
	rmc.ManageChainID = manageChainId
	rmc.KeyToSign = key3
}

func (rmc *RegisterManageChain) FactomEntry() *factom.Entry {
	extIDs := VersionAndType(rmc)
	extIDs = append(extIDs, rmc.ManageChainID.Bytes())

	sig := rmc.KeyToSign.Sign(upToSig(extIDs))
	extIDs = append(extIDs, rmc.KeyToSign.Public.Bytes())
	extIDs = append(extIDs, sig)

	e := new(factom.Entry)
	e.ExtIDs = extIDs
	e.ChainID = rmc.RootChainID.String()
	return e
}

type ManageChainMetaData struct {
	Website           *primitives.SiteURL
	LongDescription   *primitives.LongDescription
	ShortDescription  *primitives.ShortDescription
	Playlist          *common.ManyPlayList
	Thumbnail         *primitives.Image
	Banner            *primitives.Image
	ChannelTags       *primitives.TagList
	SuggestedChannels *primitives.HashList
}

func RandomManageChainMetaData() *ManageChainMetaData {
	m := NewManageChainMetaData()

	m.Website = primitives.RandomSiteURL()
	m.LongDescription = primitives.RandomLongDescription()
	m.ShortDescription = primitives.RandomShortDescription()
	m.Playlist = common.RandomManyPlayList(10)
	m.Thumbnail = primitives.RandomValidImage()
	m.Banner = primitives.RandomManyPlayList()
	m.ChannelTags = primitives.RandomTagList(uint32(constants.MAX_CHANNEL_TAGS))
	m.SuggestedChannels = primitives.RandomHashList(10)
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

// Data from Entries in MsgPack
//		Channel Website
//		Channel LongDescription
//		Channel ShortDescription
//		Playlist
//		Thumbnail
//		Banner
//		Channel Tags
//		Suggested Channels
type ManageMetaData struct {
	MetaData  ManageChainMetaData
	KeyToSign primitives.PrivateKey
	root      primitives.Hash
}

func (ManageMetaData) Type() []byte  { return TYPE_MANAGE_CHAIN_METADATA }
func (ManageMetaData) IsChain() bool { return false }
func (ManageMetaData) ForChain() int { return CHAIN_MANAGEMENT }

func (mmd *ManageMetaData) Create(metaToChange ManageChainMetaData, key3 primitives.PrivateKey, root primitives.Hash) {
	mmd.MetaData = metaToChange
	mmd.KeyToSign = key3
}

func (mmd *ManageMetaData) FactomEntry() ([]*factom.Entry, error) {
	placeHolder := make([]byte, 64)
	metaDataBytes, err := mmd.MetaData.MarshalBinary()
	if err != nil {
		return nil, err
	}

	fullContentHash := sha256.Sum256(metaDataBytes)

	tsData, err := time.Now().MarshalBinary()
	if err != nil {
		return nil, err
	}

	es := make([]*factom.Entry, 0)
	e := new(factom.Entry)
	//	0	byte		Version
	//	1	[32]byte	"Channel Management Metadata Main"
	//	2	[32]byte	RootChainID
	//	3	[32]byte	FullContentHash
	//	4	[8]byte		EntryCount
	//	5	[32]byte	ContentHash
	//	6	[15]byte	Timestamp
	//	7	[32]byte	PublicKey(3)
	//	8	[64]byte	Signature

	extIDs := VersionAndType(mmd)
	extIDs = append(extIDs, mmd.root.Bytes())
	extIDs = append(extIDs, fullContentHash[:])
	extIDs = append(extIDs, placeHolder[:8])  // EntryCount
	extIDs = append(extIDs, placeHolder[:32]) //ContentHash
	extIDs = append(extIDs, tsData)
	extIDs = append(extIDs, mmd.KeyToSign.Public.Bytes())
	extIDs = append(extIDs, placeHolder[:64]) // Sig

	headerLength := ExIDLength(extIDs)
	contentLength := len(metaDataBytes)
	totalSize := contentLength + headerLength
	contentHeaderLen := 245 + 9*2 + 2

	var entryCount int = 0
	if totalSize > constants.ENTRY_MAX_SIZE {
		entryCount = howManyEntries(headerLength, contentLength, contentHeaderLen)
	}
	extIDs[4] = primitives.Uint32ToBytes(uint32(entryCount))
	e.ExtIDs = extIDs

	if entryCount == 0 {
		e.extIDs[5] = fullContentHash[:]
		e.Content = metaDataBytes
		metaDataBytes = []byte{}
	} else {
		fl := constants.ENTRY_MAX_SIZE - ExIDLength(e.ExtIDs)
		flData := metaDataBytes[:fl]
		metaDataBytes = metaDataBytes[fl:]
		e.Content = flData
		partHash := sha256.Sum256(flData)
		e.ExtIDs[5] = partHash[:]
	}
	// Now sign full data set
	buf := new(bytes.Buffer)
	for i := 0; i < 7; i++ {
		buf.Write(e.ExtIDs[i])
	}
	contentData := buf.Next(buf.Len())
	sig = mmd.KeyToSign.Sign(contentData)
	e.ExtIDs[8] = sig

	// metaDataBytes is remaining bytes to be stiched
	bytesPerEntry := constants.ENTRY_MAX_SIZE - contentHeaderLen
	c.Entries = make([]*factom.Entry, entryCount)

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
	var seq uint32 = 0
	for len(metaDataBytes) > 0 {
		entry := new(factom.Entry)
		end := bytesPerEntry
		if len(metaDataBytes) < bytesPerEntry {
			end = len(metaDataBytes)
		}

		contentData := metaDataBytes[:end]
		metaDataBytes = metaDataBytes[end:]

		partHash := sha256.Sum256(contentData)
		// Set headers
		entry.ExtIDs = append(entry.ExtIDs, []byte{constants.FACTOM_VERSION})            // 0
		entry.ExtIDs = append(entry.ExtIDs, []byte("Channel Management Metadata Stich")) // 1
		entry.ExtIDs = append(entry.ExtIDs, root.Bytes())                                // 2
		entry.ExtIDs = append(entry.ExtIDs, fullHash[:])                                 // 3
		entry.ExtIDs = append(entry.ExtIDs, primitives.Uint32ToBytes(seq+1))             // 4 - Seq
		entry.ExtIDs = append(entry.ExtIDs, partHash[:])                                 // 5
		entry.ExtIDs = append(entry.ExtIDs, tsData)                                      // 6

		msg := upToSig(entry.ExtIDs)
		entry.ExtIDs = append(entry.ExtIDs, sigKey.Public.Bytes()) // 7
		sig := sigKey.Sign(msg)
		entry.ExtIDs = append(entry.ExtIDs, sig) // 8
		entry.ChainID = manage.String()
		entry.Content = contentData

		if int(seq) >= len(c.Entries) {
			return fmt.Errorf("Ran out of entries. Seq is %d. Entrycount is %d, %d bytes left to write", seq-1, entryCount, len(data))
		}
		c.Entries[seq] = entry
		//fmt.Println("Seq", seq, len(c.Entries[seq].Content), len(data), entry.ExtIDs[4])
		seq++
	}

	return es, nil
}

/*

// Content time


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
	c.Entries = make([]*factom.Entry, entryCount)
	bytesPerEntry := constants.ENTRY_MAX_SIZE - contentHeaderLen
	var seq uint32 = 0
	for len(data) > 0 {
		entry := new(factom.Entry)
		end := bytesPerEntry
		if len(data) < bytesPerEntry {
			end = len(data)
		}

		contentData := data[:end]
		data = data[end:]

		partHash := sha256.Sum256(contentData)
		// Set headers
		entry.ExtIDs = append(entry.ExtIDs, []byte{constants.FACTOM_VERSION})            // 0
		entry.ExtIDs = append(entry.ExtIDs, []byte("Channel Management Metadata Stich")) // 1
		entry.ExtIDs = append(entry.ExtIDs, root.Bytes())                                // 2
		entry.ExtIDs = append(entry.ExtIDs, fullHash[:])                                 // 3
		entry.ExtIDs = append(entry.ExtIDs, primitives.Uint32ToBytes(seq+1))             // 4 - Seq
		entry.ExtIDs = append(entry.ExtIDs, partHash[:])                                 // 5
		entry.ExtIDs = append(entry.ExtIDs, tsData)                                      // 6

		msg := upToSig(entry.ExtIDs)
		entry.ExtIDs = append(entry.ExtIDs, sigKey.Public.Bytes()) // 7
		sig := sigKey.Sign(msg)
		entry.ExtIDs = append(entry.ExtIDs, sig) // 8
		entry.ChainID = manage.String()
		entry.Content = contentData

		if int(seq) >= len(c.Entries) {
			return fmt.Errorf("Ran out of entries. Seq is %d. Entrycount is %d, %d bytes left to write", seq-1, entryCount, len(data))
		}
		c.Entries[seq] = entry
		//fmt.Println("Seq", seq, len(c.Entries[seq].Content), len(data), entry.ExtIDs[4])
		seq++
	}

	// Done
	c.MainEntry = e
	r.MetaData = c
*/

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

/*

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

// Data from Entries in MsgPack
//		Channel Website
//		Channel LongDescription
//		Channel ShortDescription
//		Playlist
//		Thumbnail
//		Banner
//		Channel Tags
//		Suggested Channels
func (r *ManageChain) CreateMetadata(meta *ManageChainMetaData, root primitives.Hash, manage primitives.Hash, sigKey primitives.PrivateKey) error {
	e := new(factom.Entry)

	data, err := meta.MarshalBinary()
	if err != nil {
		return err
	}

	tsData, err := time.Now().MarshalBinary()
	if err != nil {
		return err
	}

	fullHash := sha256.Sum256(data)
	contentHash := make([]byte, 32)
	holder := make([]byte, 8)

	e.ExtIDs = append(e.ExtIDs, []byte{constants.FACTOM_VERSION})           // 0
	e.ExtIDs = append(e.ExtIDs, []byte("Channel Management Metadata Main")) // 1
	e.ExtIDs = append(e.ExtIDs, root.Bytes())                               // 2
	e.ExtIDs = append(e.ExtIDs, fullHash[:])                                // 3
	e.ExtIDs = append(e.ExtIDs, holder)                                     // 4 - Holder
	e.ExtIDs = append(e.ExtIDs, contentHash)                                // 5 - Holder
	e.ExtIDs = append(e.ExtIDs, tsData)                                     // 6

	msg := upToSig(e.ExtIDs)
	e.ExtIDs = append(e.ExtIDs, sigKey.Public.Bytes()) // 7
	sig := sigKey.Sign(msg)
	e.ExtIDs = append(e.ExtIDs, sig) // 8

	// Content time
	headerLength := ExIDLength(e.ExtIDs)
	contentLength := len(data)
	totalSize := contentLength + headerLength
	contentHeaderLen := 245 + 9*2 + 2

	var entryCount int = 0
	if totalSize > constants.ENTRY_MAX_SIZE {
		entryCount = howManyEntries(headerLength, contentLength, contentHeaderLen)
	}
	e.ExtIDs[4] = primitives.Uint32ToBytes(uint32(entryCount))
	e.ChainID = manage.String()

	c := new(ChanMetaDataEntries)
	if entryCount == 0 {
		e.ExtIDs[5] = fullHash[:]
		e.Content = data
		data = []byte{}
	} else {
		fl := constants.ENTRY_MAX_SIZE - ExIDLength(e.ExtIDs)
		flData := data[:fl]
		data = data[fl:]
		e.Content = flData
		partHash := sha256.Sum256(flData)
		e.ExtIDs[5] = partHash[:]
	}

	// Redo signature with  new values
	buf := new(bytes.Buffer)
	for i := 0; i < 7; i++ {
		buf.Write(e.ExtIDs[i])
	}
	contentData := buf.Next(buf.Len())
	sig = sigKey.Sign(contentData)
	e.ExtIDs[8] = sig

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
	c.Entries = make([]*factom.Entry, entryCount)
	bytesPerEntry := constants.ENTRY_MAX_SIZE - contentHeaderLen
	var seq uint32 = 0
	for len(data) > 0 {
		entry := new(factom.Entry)
		end := bytesPerEntry
		if len(data) < bytesPerEntry {
			end = len(data)
		}

		contentData := data[:end]
		data = data[end:]

		partHash := sha256.Sum256(contentData)
		// Set headers
		entry.ExtIDs = append(entry.ExtIDs, []byte{constants.FACTOM_VERSION})            // 0
		entry.ExtIDs = append(entry.ExtIDs, []byte("Channel Management Metadata Stich")) // 1
		entry.ExtIDs = append(entry.ExtIDs, root.Bytes())                                // 2
		entry.ExtIDs = append(entry.ExtIDs, fullHash[:])                                 // 3
		entry.ExtIDs = append(entry.ExtIDs, primitives.Uint32ToBytes(seq+1))             // 4 - Seq
		entry.ExtIDs = append(entry.ExtIDs, partHash[:])                                 // 5
		entry.ExtIDs = append(entry.ExtIDs, tsData)                                      // 6

		msg := upToSig(entry.ExtIDs)
		entry.ExtIDs = append(entry.ExtIDs, sigKey.Public.Bytes()) // 7
		sig := sigKey.Sign(msg)
		entry.ExtIDs = append(entry.ExtIDs, sig) // 8
		entry.ChainID = manage.String()
		entry.Content = contentData

		if int(seq) >= len(c.Entries) {
			return fmt.Errorf("Ran out of entries. Seq is %d. Entrycount is %d, %d bytes left to write", seq-1, entryCount, len(data))
		}
		c.Entries[seq] = entry
		//fmt.Println("Seq", seq, len(c.Entries[seq].Content), len(data), entry.ExtIDs[4])
		seq++
	}

	// Done
	c.MainEntry = e
	r.MetaData = c
	return nil
}

type ManageChainMetaData struct {
	Website           *primitives.SiteURL
	LongDescription   *primitives.LongDescription
	ShortDescription  *primitives.ShortDescription
	Playlist          *common.ManyPlayList
	Thumbnail         *primitives.Image
	Banner            *primitives.Image
	ChannelTags       *primitives.TagList
	SuggestedChannels *primitives.HashList
}

func NewManageChainMetaData() *ManageChainMetaData {
	m := new(ManageChainMetaData)
	m.Website = new(primitives.SiteURL)
	m.LongDescription = new(primitives.LongDescription)
	m.ShortDescription = new(primitives.ShortDescription)
	m.Playlist = new(common.ManyPlayList)
	m.Thumbnail = new(primitives.Image)
	m.Banner = new(primitives.Image)
	m.ChannelTags = new(primitives.TagList)
	m.SuggestedChannels = new(primitives.HashList)

	return m
}

func RandomManageChainMetaData() *ManageChainMetaData {
	m := NewManageChainMetaData()

	m.Website = primitives.RandomSiteURL()
	m.LongDescription = primitives.RandomLongDescription()
	m.ShortDescription = primitives.RandomShortDescription()
	m.Playlist = common.RandomManyPlayList(10)
	m.Thumbnail = primitives.RandomImage()
	m.Banner = primitives.RandomImage()
	m.ChannelTags = primitives.RandomTagList(uint32(constants.MAX_CHANNEL_TAGS))
	m.SuggestedChannels = primitives.RandomHashList(10)
	return m
}

func RandomHugeManageChainMetaData() *ManageChainMetaData {
	m := NewManageChainMetaData()

	m.Website = primitives.RandomSiteURL()
	m.LongDescription = primitives.RandomLongDescription()
	m.ShortDescription = primitives.RandomShortDescription()
	m.Playlist = common.RandomManyPlayList(10)
	m.Thumbnail = primitives.RandomHugeImage()
	m.Banner = primitives.RandomHugeImage()
	m.ChannelTags = primitives.RandomTagList(uint32(constants.MAX_CHANNEL_TAGS))
	m.SuggestedChannels = primitives.RandomHashList(10)
	return m
}

func encodeBytes(data []byte) []byte {
	return data
	// return string(data)
}

func (m *ManageChainMetaData) MarshalBinary() ([]byte, error) {
	mb := new(ManageChainMetaDataBytes)
	data, err := m.Website.MarshalBinary()
	if err != nil {
		return nil, err
	}
	mb.Website = encodeBytes(data)

	data, err = m.LongDescription.MarshalBinary()
	if err != nil {
		return nil, err
	}
	mb.LongDescription = encodeBytes(data)

	data, err = m.ShortDescription.MarshalBinary()
	if err != nil {
		return nil, err
	}
	mb.ShortDescription = encodeBytes(data)

	data, err = m.Playlist.MarshalBinary()
	if err != nil {
		return nil, err
	}
	mb.Playlist = encodeBytes(data)

	data, err = m.Thumbnail.MarshalBinary()
	if err != nil {
		return nil, err
	}
	mb.Thumbnail = encodeBytes(data)

	data, err = m.Banner.MarshalBinary()
	if err != nil {
		return nil, err
	}
	mb.Banner = encodeBytes(data)

	data, err = m.ChannelTags.MarshalBinary()
	if err != nil {
		return nil, err
	}
	mb.ChannelTags = encodeBytes(data)

	data, err = m.SuggestedChannels.MarshalBinary()
	if err != nil {
		return nil, err
	}
	mb.SuggestedChannels = encodeBytes(data)

	msgPackData, err := mb.MarshalMsg(nil)
	if err != nil {
		return nil, err
	}

	length := primitives.Uint32ToBytes(uint32(len(msgPackData)))
	buf := new(bytes.Buffer)
	buf.Write(length)
	buf.Write(msgPackData)

	return buf.Next(buf.Len()), nil
}

func (m *ManageChainMetaData) UnmarshalBinary(data []byte) (err error) {
	_, err = m.UnmarshalBinaryData(data)
	return
}

func (m *ManageChainMetaData) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	mb := new(ManageChainMetaDataBytes)
	newData = data

	u, err := primitives.BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}
	newData = newData[4:]

	_, err = mb.UnmarshalMsg(newData[:u])
	if err != nil {
		return data, err
	}

	newData = newData[u:]

	if len(mb.Website) > 0 {
		m.Website = new(primitives.SiteURL)
		_, err = m.Website.UnmarshalBinaryData(mb.Website)
		if err != nil {
			return data, err
		}
	} else {
		m.Website = nil
	}

	if len(mb.LongDescription) > 0 {
		m.LongDescription = new(primitives.LongDescription)
		err = m.LongDescription.UnmarshalBinary(mb.LongDescription)
		if err != nil {
			return data, err
		}
	} else {
		m.LongDescription = nil
	}

	if len(mb.ShortDescription) > 0 {
		m.ShortDescription = new(primitives.ShortDescription)
		err = m.ShortDescription.UnmarshalBinary(mb.ShortDescription)
		if err != nil {
			return data, err
		}
	} else {
		m.ShortDescription = nil
	}

	if len(mb.Playlist) > 0 {
		m.Playlist = new(common.ManyPlayList)
		err = m.Playlist.UnmarshalBinary(mb.Playlist)
		if err != nil {
			return data, err
		}
	} else {
		m.Playlist = nil
	}

	if len(mb.Thumbnail) > 0 {
		m.Thumbnail = new(primitives.Image)
		err = m.Thumbnail.UnmarshalBinary(mb.Thumbnail)
		if err != nil {
			return data, err
		}
	} else {
		m.Thumbnail = nil
	}

	if len(mb.Banner) > 0 {
		m.Banner = new(primitives.Image)
		err = m.Banner.UnmarshalBinary(mb.Banner)
		if err != nil {
			return data, err
		}
	} else {
		m.Banner = nil
	}

	if len(mb.ChannelTags) > 0 {
		m.ChannelTags = new(primitives.TagList)
		err = m.ChannelTags.UnmarshalBinary(mb.ChannelTags)
		if err != nil {
			return data, err
		}
	} else {
		m.ChannelTags = nil
	}

	if len(mb.SuggestedChannels) > 0 {
		m.SuggestedChannels = new(primitives.HashList)
		err = m.SuggestedChannels.UnmarshalBinary(mb.SuggestedChannels)
		if err != nil {
			return data, err
		}
	} else {
		m.SuggestedChannels = nil
	}

	return
}

// nilComp returns:
//		0 	Both nil		Skip
//		1 	1 nil			Return false
//		2 	none nil		Compare
func nilComp(a interface{}, b interface{}) int {
	if isNil(a) && isNil(b) {
		return 0
	}
	if !isNil(a) && !isNil(b) {
		return 2
	}
	return 1
}

func isNil(o interface{}) bool {
	if !reflect.ValueOf(o).Elem().IsValid() {
		return true
	}
	return false
}

func (a *ManageChainMetaData) IsSameAs(b *ManageChainMetaData) bool {
	if nilComp(a.Website, b.Website) != 0 &&
		(nilComp(a.Website, b.Website) == 1 || !a.Website.IsSameAs(b.Website)) {
		return false
	}

	if nilComp(a.LongDescription, b.LongDescription) != 0 &&
		(nilComp(a.LongDescription, b.LongDescription) == 1 || !a.LongDescription.IsSameAs(b.LongDescription)) {
		return false
	}

	if nilComp(a.ShortDescription, b.ShortDescription) != 0 &&
		(nilComp(a.ShortDescription, b.ShortDescription) == 1 || !a.ShortDescription.IsSameAs(b.ShortDescription)) {
		return false
	}

	if nilComp(a.Playlist, b.Playlist) != 0 &&
		(nilComp(a.Playlist, b.Playlist) == 1 || !a.Playlist.IsSameAs(b.Playlist)) {
		return false
	}

	if nilComp(a.Thumbnail, b.Thumbnail) != 0 &&
		(nilComp(a.Thumbnail, b.Thumbnail) == 1 || !a.Thumbnail.IsSameAs(b.Thumbnail)) {
		return false
	}

	if nilComp(a.Banner, b.Banner) != 0 &&
		(nilComp(a.Banner, b.Banner) == 1 || !a.Banner.IsSameAs(b.Banner)) {
		return false
	}

	if nilComp(a.ChannelTags, b.ChannelTags) != 0 &&
		(nilComp(a.ChannelTags, b.ChannelTags) == 1 || !a.ChannelTags.IsSameAs(b.ChannelTags)) {
		return false
	}

	if nilComp(a.SuggestedChannels, b.SuggestedChannels) != 0 &&
		(nilComp(a.SuggestedChannels, b.SuggestedChannels) == 1 || !a.SuggestedChannels.IsSameAs(b.SuggestedChannels)) {
		return false
	}

	return true
}
*/
