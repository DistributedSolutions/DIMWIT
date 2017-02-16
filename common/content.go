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
	Type             byte
	ContentID        primitives.Hash
	RootChainID      primitives.Hash
	ContentTitle     primitives.Title
	InfoHash         primitives.InfoHash
	LongDescription  primitives.LongDescription
	ShortDescription primitives.ShortDescription
	ActionFiles      primitives.FileList
	Thumbnail        primitives.Image
	Series           byte
	// There can be lots of parts
	Part [2]byte
	Tags primitives.TagList

	// Retrieved from Blockchain
	CreationTime time.Time

	// Torrent Metadata
	Trackers primitives.TrackerList
	FileList primitives.FileList
}

func (cl *ContentList) GetContents() []Content {
	return cl.contentList
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

func (a *Content) IsSameAs(b *Content) bool {
	if a.Type != b.Type {
		return false
	}

	if !a.ContentID.IsSameAs(&b.ContentID) {
		return false
	}

	if !a.RootChainID.IsSameAs(&b.RootChainID) {
		return false
	}

	if !a.ContentTitle.IsSameAs(&b.ContentTitle) {
		return false
	}

	if !a.InfoHash.IsSameAs(&b.InfoHash) {
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

	if !a.Trackers.IsSameAs(&b.Trackers) {
		return false
	}

	if !a.FileList.IsSameAs(&b.FileList) {
		return false
	}

	if a.CreationTime.Nanosecond() != b.CreationTime.Nanosecond() {
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
	length      uint32
	contentList []Content
}

func RandomContentList(max uint32) *ContentList {
	p := new(ContentList)
	l := random.RandomUInt32Between(0, max)
	p.length = l

	p.contentList = make([]Content, l)
	for i := range p.contentList {
		p.contentList[i] = *RandomNewContent()
	}

	return p
}

func (a *ContentList) IsSameAs(b *ContentList) bool {
	if a.length != b.length {
		return false
	}

	for i := range a.contentList {
		if !a.contentList[i].IsSameAs(&b.contentList[i]) {
			return false
		}
	}
	return true
}

func (p *ContentList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data := primitives.Uint32ToBytes(p.length)
	buf.Write(data)

	for i := range p.contentList {
		data, err := p.contentList[i].MarshalBinary()
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
	p.length = u
	newData = newData[4:]

	p.contentList = make([]Content, u)
	var i uint32
	for i = 0; i < u; i++ {
		newData, err = p.contentList[i].UnmarshalBinaryData(newData)
		if err != nil {
			return data, err
		}
	}

	return
}
