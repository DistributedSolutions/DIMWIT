package elements

import (
	"bytes"
	"crypto/sha256"
	"fmt"
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

func (m *Manage) FactomElements() ([]*factom.Entry, *factom.Chain) {
	es := make([]*factom.Entry, 0)
	c := m.ManageChain.FactomChain()

	es = append(es, m.RegisterManageChain.FactomEntry())

	return es, c
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
	m.KeyToSign = key3

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
	m.Thumbnail = primitives.RandomValidImage(constants.MAX_IMAGE_SIZE)
	m.Banner = primitives.RandomValidImage(constants.MAX_IMAGE_SIZE)
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
	MetaData  *ManageChainMetaData
	KeyToSign primitives.PrivateKey
	root      primitives.Hash
	manage    primitives.Hash
}

func (ManageMetaData) Type() []byte  { return TYPE_MANAGE_CHAIN_METADATA }
func (ManageMetaData) IsChain() bool { return false }
func (ManageMetaData) ForChain() int { return CHAIN_MANAGEMENT }

func (mmd *ManageMetaData) Create(metaToChange *ManageChainMetaData, key3 primitives.PrivateKey, root primitives.Hash, manChain primitives.Hash) {
	mmd.MetaData = metaToChange
	mmd.KeyToSign = key3
	mmd.manage = manChain
	mmd.root = root
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
		e.ExtIDs[5] = fullContentHash[:]
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
	sig := mmd.KeyToSign.Sign(contentData)
	e.ExtIDs[8] = sig
	e.ChainID = mmd.manage.String()

	es = append(es, e)
	// metaDataBytes is remaining bytes to be stiched
	bytesPerEntry := constants.ENTRY_MAX_SIZE - contentHeaderLen
	stiches := make([]*factom.Entry, entryCount)

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
		entry.ExtIDs = append(entry.ExtIDs, mmd.root.Bytes())                            // 2
		entry.ExtIDs = append(entry.ExtIDs, fullContentHash[:])                          // 3
		entry.ExtIDs = append(entry.ExtIDs, primitives.Uint32ToBytes(seq+1))             // 4 - Seq
		entry.ExtIDs = append(entry.ExtIDs, partHash[:])                                 // 5
		entry.ExtIDs = append(entry.ExtIDs, tsData)                                      // 6

		msg := upToSig(entry.ExtIDs)
		entry.ExtIDs = append(entry.ExtIDs, mmd.KeyToSign.Public.Bytes()) // 7
		sig := mmd.KeyToSign.Sign(msg)
		entry.ExtIDs = append(entry.ExtIDs, sig) // 8
		entry.ChainID = mmd.manage.String()
		entry.Content = contentData

		if int(seq) >= len(stiches) {
			return nil, fmt.Errorf("Ran out of entries. Seq is %d. Entrycount is %d, %d bytes left to write", seq-1, entryCount, len(metaDataBytes))
		}
		stiches[seq] = entry
		//fmt.Println("Seq", seq, len(c.Entries[seq].Content), len(data), entry.ExtIDs[4])
		seq++
	}

	es = append(es, stiches...)
	return es, nil
}
