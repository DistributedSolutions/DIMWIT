package common

import (
	"bytes"
	"fmt"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type Channel struct {
	RootChainID       primitives.Hash
	ManagementChainID primitives.Hash
	ContentChainID    primitives.Hash
	// They are not an array, because they are never referenced as an array
	LV1PublicKey      primitives.PublicKey // Critical
	LV2PublicKey      primitives.PublicKey // Critical
	LV3PublicKey      primitives.PublicKey // Critical
	ContentSingingKey primitives.PublicKey // Critical
	ChannelTitle      primitives.Title     // Critical

	Website          primitives.SiteURL          // Not-Critical
	LongDescription  primitives.LongDescription  // Not-Critical
	ShortDescription primitives.ShortDescription // Not-Critical
	Playlist         ManyPlayList                // Not-Critical
	Thumbnail        primitives.Image            // Not-Critical
	Banner           primitives.Image            // Not-Critical
	Tags             primitives.TagList          // Not-Critical
	SuggestedChannel primitives.HashList         // Not-Critical
	Content          ContentList                 // Not-Critical

	CreationTime time.Time // Not-Critical
}

func NewChannel() *Channel {
	c := new(Channel)

	c.Tags = *primitives.NewTagList(uint32(constants.MAX_CHANNEL_TAGS))
	return c
}

func RandomNewChannel() *Channel {
	c := new(Channel)

	c.RootChainID = *primitives.RandomHash()
	c.ManagementChainID = *primitives.RandomHash()
	c.ContentChainID = *primitives.RandomHash()
	c.LV1PublicKey = *primitives.RandomPublicKey()
	c.LV2PublicKey = *primitives.RandomPublicKey()
	c.LV3PublicKey = *primitives.RandomPublicKey()
	c.ContentSingingKey = *primitives.RandomPublicKey()

	c.ChannelTitle = *primitives.RandomTitle()
	c.Website = *primitives.RandomSiteURL()
	c.LongDescription = *primitives.RandomLongDescription()
	c.ShortDescription = *primitives.RandomShortDescription()
	c.Content = *SmartRandomContentList(random.RandomUInt32Between(0, 30),
		c.RootChainID,
		c.ContentChainID)
	// c.Playlist = *RandomManyPlayList(random.RandomUInt32Between(0, 100))
	c.Playlist = *SmartRandomManyPlayList(random.RandomUInt32Between(0, 100), c.Content)
	c.Thumbnail = *primitives.RandomImage()
	c.Banner = *primitives.RandomImage()
	c.Tags = *primitives.RandomTagList(uint32(constants.MAX_CHANNEL_TAGS))
	c.SuggestedChannel = *primitives.RandomHashList(random.RandomUInt32Between(0, 100))
	c.CreationTime = time.Now()

	return c
}

func (a *Channel) Status() int {
	if a.full() {
		return constants.CHANNEL_FULL
	}

	if a.ready() {
		return constants.CHANNEL_READY
	}

	return constants.CHANNEL_NOT_READY
}

// Woo! full!
func (a *Channel) full() bool {
	if !a.ready() {
		return false
	}

	if a.Website.Empty() {
		return false
	}

	if a.LongDescription.Empty() {
		return false
	}

	if a.ShortDescription.Empty() {
		return false
	}

	if a.Thumbnail.Empty() {
		return false
	}

	if a.Banner.Empty() {
		return false
	}

	if a.Tags.Empty() {
		return false
	}

	return true
}

// Reaady for public consumption
func (a *Channel) ready() bool {
	if a.LV1PublicKey.Empty() {
		return false
	}

	if a.LV1PublicKey.Empty() {
		return false
	}

	if a.LV2PublicKey.Empty() {
		return false
	}

	if a.LV3PublicKey.Empty() {
		return false
	}

	if a.ContentSingingKey.Empty() {
		return false
	}

	if a.ChannelTitle.Empty() {
		return false
	}

	return true
}

func (a *Channel) IsSameAs(b *Channel) bool {
	if !a.RootChainID.IsSameAs(&b.RootChainID) {
		fmt.Println("Exit 1")
		return false
	}

	if !a.ManagementChainID.IsSameAs(&b.ManagementChainID) {
		fmt.Println("Exit 1")
		return false
	}

	if !a.ContentChainID.IsSameAs(&b.ContentChainID) {
		fmt.Println("Exit 2")
		return false
	}

	if !a.LV1PublicKey.IsSameAs(&b.LV1PublicKey) {
		fmt.Println("Exit 3")
		return false
	}

	if !a.LV2PublicKey.IsSameAs(&b.LV2PublicKey) {
		fmt.Println("Exit 4")
		return false
	}

	if !a.LV3PublicKey.IsSameAs(&b.LV3PublicKey) {
		fmt.Println("Exit 5")
		return false
	}

	if !a.ContentSingingKey.IsSameAs(&b.ContentSingingKey) {
		fmt.Println("Exit 6")
		return false
	}

	if !a.Website.IsSameAs(&b.Website) {
		fmt.Printf("Exit 7. Web A: %s, Web B: %s\n", a.Website.String(), b.Website.String())
		return false
	}

	if !a.LongDescription.IsSameAs(&b.LongDescription) {
		fmt.Println("Exit 8")
		return false
	}

	if !a.ShortDescription.IsSameAs(&b.ShortDescription) {
		fmt.Println("Exit 9")
		return false
	}

	if !a.Playlist.IsSameAs(&b.Playlist) {
		fmt.Println("Exit 10")
		return false
	}

	if !a.Thumbnail.IsSameAs(&b.Thumbnail) {
		fmt.Println("Exit 11")
		return false
	}

	if !a.Banner.IsSameAs(&b.Banner) {
		fmt.Println("Exit 12")
		return false
	}

	if !a.Tags.IsSameAs(&b.Tags) {
		fmt.Println(a.Tags, b.Tags)
		fmt.Println("Exit 13")
		return false
	}

	if !a.SuggestedChannel.IsSameAs(&b.SuggestedChannel) {
		fmt.Println("Exit 14")
		return false
	}

	if !a.Content.IsSameAs(&b.Content) {
		fmt.Println("Exit 15")
		fmt.Println("DEBUG:", len(a.Content.GetContents()), len(b.Content.ContentList))
		return false
	}

	/*if a.CreationTime.Unix() != b.CreationTime.Unix() {
		fmt.Println("Exit 16")
		return false
	}*/

	return true
}

func (c *Channel) MarshalBinary() (data []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[Channel] A panic has occurred while marshaling: %s", r)
			return
		}
	}()

	buf := new(bytes.Buffer)

	data, err = c.RootChainID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.ManagementChainID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.ContentChainID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.LV1PublicKey.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.LV2PublicKey.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.LV3PublicKey.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.ContentSingingKey.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.ChannelTitle.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.Website.MarshalBinary()
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

	data, err = c.Playlist.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.Thumbnail.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.Banner.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.Tags.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.SuggestedChannel.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.Content.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = c.CreationTime.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	return buf.Next(buf.Len()), nil
}

func (c *Channel) UnmarshalBinary(data []byte) (err error) {
	_, err = c.UnmarshalBinaryData(data)
	return err
}

func (c *Channel) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[channel] A panic has occurred while marshaling: %s", r)
			return
		}
	}()

	newData = data

	newData, err = c.RootChainID.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.ManagementChainID.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.ContentChainID.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.LV1PublicKey.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.LV2PublicKey.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.LV3PublicKey.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.ContentSingingKey.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.ChannelTitle.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.Website.UnmarshalBinaryData(newData)
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

	newData, err = c.Playlist.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.Thumbnail.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.Banner.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.Tags.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.SuggestedChannel.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = c.Content.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	var tmp time.Time
	err = tmp.UnmarshalBinary(newData[:15])
	if err != nil {
		fmt.Printf("Hit the error here: %x\n", newData[:15])
		return data, err
	}
	c.CreationTime = tmp
	newData = newData[15:]

	return newData, nil
}
