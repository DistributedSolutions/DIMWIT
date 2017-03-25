package common

import (
	"bytes"
	"fmt"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type Content struct {
	// Content Metadada
	Type             byte                        `json:"type"`
	ContentID        primitives.Hash             `json:"contentchain"`
	RootChainID      primitives.Hash             `json:"rootchain"`
	ContentTitle     primitives.Title            `json:"title"`
	InfoHash         primitives.InfoHash         `json:"infohash"`
	LongDescription  primitives.LongDescription  `json:"longdesc"`
	ShortDescription primitives.ShortDescription `json:"shortdesc"`
	ActionFiles      primitives.FileList         `json:"actionfiles"`
	Thumbnail        primitives.Image            `json:"thumbnail"`
	Series           byte                        `json:"series"`
	// There can be lots of parts
	Part [2]byte            `json:"part"`
	Tags primitives.TagList `json:"tags"`

	// Retrieved from Blockchain
	CreationTime time.Time `json:"creationtime"`

	// Torrent Metadata
	Trackers primitives.TrackerList `json:"trackerlist"`
	FileList primitives.FileList    `json:"filelist"`
}

func (cl *ContentList) GetContents() []Content {
	return cl.ContentList
}

func NewContent() *Content {
	c := new(Content)
	c.Tags = *primitives.RandomTagList(uint32(constants.MAX_CONTENT_TAGS))
	return c

}

func RandomNewContent() *Content {
	c := new(Content)

	c.ContentID = *primitives.RandomHash()
	c.RootChainID = *primitives.RandomHash()
	c.ContentTitle = *primitives.RandomTitle()
	c.InfoHash = *primitives.RandomInfoHash()
	c.LongDescription = *primitives.RandomLongDescription()
	c.ShortDescription = *primitives.RandomShortDescription()
	c.ActionFiles = *primitives.RandomFileList(uint32(10))
	c.Thumbnail = *primitives.RandomImage()
	c.Tags = *primitives.RandomTagList(uint32(constants.MAX_CONTENT_TAGS))
	c.Trackers = *primitives.RandomTrackerList(uint32(5))
	c.FileList = *primitives.RandomFileList(uint32(10))
	c.CreationTime = time.Now()

	return c
}

func SmartRandomNewContent(root primitives.Hash, content primitives.Hash) *Content {
	c := new(Content)

	c.ContentID = content
	c.RootChainID = root
	c.ContentTitle = *primitives.RandomTitle()
	c.InfoHash = *primitives.RandomInfoHash()
	c.LongDescription = *primitives.RandomLongDescription()
	c.ShortDescription = *primitives.RandomShortDescription()
	c.ActionFiles = *primitives.RandomFileList(uint32(10))
	c.Thumbnail = *primitives.RandomHugeImage()
	c.Tags = *primitives.RandomTagList(uint32(constants.MAX_CONTENT_TAGS))
	c.Trackers = *primitives.RandomTrackerList(uint32(5))
	c.FileList = *primitives.RandomFileList(uint32(10))
	c.CreationTime = time.Now()

	return c
}

func (a *Content) IsSameAs(b *Content) bool {
	if (a == nil && b != nil) || (b == nil && a != nil) {
		return false
	}

	if a.Type != b.Type {
		//fmt.Println(a.Type, b.Type)
		//log.Println("Content IsSameAs Exit 1")
		return false
	}

	if !a.ContentID.IsSameAs(&b.ContentID) {
		//fmt.Println(a.ContentID.String(), b.ContentID.String())
		//log.Println("Content IsSameAs Exit 2")
		return false
	}

	if !a.RootChainID.IsSameAs(&b.RootChainID) {
		//log.Println("Content IsSameAs Exit 3")
		//fmt.Println(a.RootChainID.String(), b.RootChainID.String())
		return false
	}

	if !a.ContentTitle.IsSameAs(&b.ContentTitle) {
		//log.Println("Content IsSameAs Exit 4")
		return false
	}

	if !a.InfoHash.IsSameAs(&b.InfoHash) {
		//log.Println("Content IsSameAs Exit 5")
		return false
	}

	if !a.LongDescription.IsSameAs(&b.LongDescription) {
		//log.Println("Content IsSameAs Exit 6")
		return false
	}

	if !a.ShortDescription.IsSameAs(&b.ShortDescription) {
		//log.Println("Content IsSameAs Exit 7")
		return false
	}

	if !a.ActionFiles.IsSameAs(&b.ActionFiles) {
		//log.Println("Content IsSameAs Exit 8")
		return false
	}

	if !a.Thumbnail.IsSameAs(&b.Thumbnail) {
		//log.Println("Content IsSameAs Exit 9")
		return false
	}

	if a.Series != b.Series {
		//log.Println("Content IsSameAs Exit 10")
		return false
	}

	if a.Part[0] != b.Part[0] || a.Part[1] != b.Part[1] {
		//log.Println("Content IsSameAs Exit 11")
		return false
	}

	if !a.Tags.IsSameAs(&b.Tags) {
		//log.Println("Content IsSameAs Exit 12")
		return false
	}

	if !a.Trackers.IsSameAs(&b.Trackers) {
		//log.Println("Content IsSameAs Exit 13")
		return false
	}

	if !a.FileList.IsSameAs(&b.FileList) {
		//log.Println("Content IsSameAs Exit 14")
		return false
	}

	diff := a.CreationTime.Unix() - b.CreationTime.Unix()
	if diff < 0 {
		diff = -1 * diff
	}
	if diff > 60*60*24 { // 1 day difference
		//log.Println("Content IsSameAs Exit 15")
		return false
	}

	return true
}

func (c *Content) MarshalBinary() (data []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[Content] A panic has occurred while marshaling: %s", r)
			return
		}
	}()

	buf := new(bytes.Buffer)

	buf.Write([]byte{c.Type})

	data, err = c.ContentID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.RootChainID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.ContentTitle.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.InfoHash.MarshalBinary()
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

	data, err = c.Trackers.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.FileList.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.CreationTime.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	return buf.Next(buf.Len()), err
}

func (c *Content) UnmarshalBinary(data []byte) error {
	_, err := c.UnmarshalBinaryData(data)
	return err
}

func (c *Content) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data

	c.Type = newData[0]
	newData = newData[1:]

	// c.ContentID = new(primitives.Hash)
	newData, err = c.ContentID.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	// c.RootChainID = new(primitives.Hash)
	newData, err = c.RootChainID.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	// c.ContentTitle = new(primitives.Hash)
	newData, err = c.ContentTitle.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	// c.InfoHash = new(primitives.InfoHash)
	newData, err = c.InfoHash.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	// c.LongDescription = new(primitives.LongDescription)
	newData, err = c.LongDescription.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	// c.ShortDescription = new(primitives.ShortDescription)
	newData, err = c.ShortDescription.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	// c.ActionFiles = new(primitives.FileList)
	newData, err = c.ActionFiles.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	// c.Thumbnail = new(primitives.Image)
	newData, err = c.Thumbnail.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	c.Series = newData[0]
	newData = newData[1:]

	copy(c.Part[:], newData[:2])
	newData = newData[2:]

	// c.Tags = new(primitives.TagList)
	newData, err = c.Tags.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	// c.Trackers = new(primitives.TrackerList)
	newData, err = c.Trackers.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	// c.FileList = new(primitives.FileList)
	newData, err = c.FileList.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	err = c.CreationTime.UnmarshalBinary(newData[:15])
	if err != nil {
		return data, err
	}
	newData = newData[15:]

	return
}

//
// List of Content
//

type ContentList struct {
	ContentList []Content `json:"contentlist"`
}

func RandomContentList(max uint32) *ContentList {
	p := new(ContentList)
	l := random.RandomUInt32Between(0, max)

	p.ContentList = make([]Content, l)
	for i := range p.ContentList {
		p.ContentList[i] = *RandomNewContent()
	}

	return p
}

func SmartRandomContentList(max uint32, root primitives.Hash, content primitives.Hash) *ContentList {
	p := new(ContentList)
	l := random.RandomUInt32Between(0, max)

	p.ContentList = make([]Content, l)
	for i := range p.ContentList {
		p.ContentList[i] = *SmartRandomNewContent(root, content)
	}

	return p
}

func (a *ContentList) IsSameAs(b *ContentList) bool {
	if len(a.ContentList) != len(b.ContentList) {
		return false
	}

	for i := range a.ContentList {
		if !a.ContentList[i].IsSameAs(&b.ContentList[i]) {
			return false
		}
	}
	return true
}

func (p *ContentList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data := primitives.Uint32ToBytes(uint32(len(p.ContentList)))
	buf.Write(data)

	for i := range p.ContentList {
		data, err := p.ContentList[i].MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(data)
	}

	return buf.Next(buf.Len()), nil
}

func (p *ContentList) UnmarshalBinary(data []byte) error {
	_, err := p.UnmarshalBinaryData(data)
	return err
}

func (p *ContentList) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data
	u, err := primitives.BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}
	newData = newData[4:]

	p.ContentList = make([]Content, u)
	var i uint32
	for i = 0; i < u; i++ {
		newData, err = p.ContentList[i].UnmarshalBinaryData(newData)
		if err != nil {
			return data, err
		}
	}

	return
}
