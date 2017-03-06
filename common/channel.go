package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = log.Prefix()

type ChannelList struct {
	List []Channel `json:"channellist"`
}

type Channel struct {
	RootChainID       primitives.Hash `json:"rootchain"`
	ManagementChainID primitives.Hash `json:"managechain"`
	ContentChainID    primitives.Hash `json:"contentchain"`
	// They are not an array, because they are never referenced as an array
	LV1PublicKey      primitives.PublicKey `json:"pubkey1"`    // Critical
	LV2PublicKey      primitives.PublicKey `json:"pubkey2"`    // Critical
	LV3PublicKey      primitives.PublicKey `json:"pubkey3"`    // Critical
	ContentSingingKey primitives.PublicKey `json:"contentkey"` // Critical
	ChannelTitle      primitives.Title     `json:"title"`      // Critical

	Website          primitives.SiteURL          `json:"site"`              // Not-Critical
	LongDescription  primitives.LongDescription  `json:"longdesc"`          // Not-Critical
	ShortDescription primitives.ShortDescription `json:"shortdesc"`         // Not-Critical
	Playlist         ManyPlayList                `json:"playlist"`          // Not-Critical
	Thumbnail        primitives.Image            `json:"thumbnail"`         // Not-Critical
	Banner           primitives.Image            `json:"banner"`            // Not-Critical
	Tags             primitives.TagList          `json:"tags"`              // Not-Critical
	SuggestedChannel primitives.HashList         `json:"suggestedchannels"` // Not-Critical
	Content          ContentList                 `json:"contentlist"`       // Not-Critical

	CreationTime time.Time `json:"creationtime"` // Not-Critical
}

func NewChannel() *Channel {
	c := new(Channel)

	c.Tags = *primitives.NewTagList(uint32(constants.MAX_CHANNEL_TAGS))
	return c
}

func RandomNewSmallChannel() *Channel {
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
	c.Content = *SmartRandomContentList(random.RandomUInt32Between(1, 3),
		c.RootChainID,
		c.ContentChainID)
	// c.Playlist = *RandomManyPlayList(random.RandomUInt32Between(0, 100))
	c.Playlist = *SmartRandomManyPlayList(random.RandomUInt32Between(1, 2), c.Content)
	c.Thumbnail = *primitives.RandomImage()
	c.Thumbnail.SetImage([]byte{0x00, 0x01, 0x02})
	c.Banner = *primitives.RandomImage()
	c.Banner.SetImage([]byte{0x00, 0x01, 0x02})
	c.Tags = *primitives.RandomTagList(uint32(constants.MAX_CHANNEL_TAGS))
	c.SuggestedChannel = *primitives.RandomHashList(random.RandomUInt32Between(1, 2))
	c.CreationTime = time.Now()

	for i, con := range c.Content.ContentList {
		con.Thumbnail.SetImage([]byte{0xFF})
		c.Content.ContentList[i] = con
	}

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
	c.Thumbnail = *primitives.RandomHugeImage()
	c.Banner = *primitives.RandomImage()
	c.Tags = *primitives.RandomTagList(uint32(constants.MAX_CHANNEL_TAGS))
	c.SuggestedChannel = *primitives.RandomHashList(random.RandomUInt32Between(0, 100))
	c.CreationTime = time.Now()

	return c
}

type CustomJSONMarshalChannel struct {
	RootChainID       string `json:"rootchain"`
	ManagementChainID string `json:"managechain"`
	ContentChainID    string `json:"contentchain"`
	// They are not an array, because they are never referenced as an array
	LV1PublicKey      string            `json:"pubkey1"`    // Critical
	LV2PublicKey      string            `json:"pubkey2"`    // Critical
	LV3PublicKey      string            `json:"pubkey3"`    // Critical
	ContentSingingKey string            `json:"contentkey"` // Critical
	ChannelTitle      *primitives.Title `json:"title"`      // Critical

	Website          *primitives.SiteURL             `json:"site"`              // Not-Critical
	LongDescription  *primitives.LongDescription     `json:"longdesc"`          // Not-Critical
	ShortDescription *primitives.ShortDescription    `json:"shortdesc"`         // Not-Critical
	Playlist         *ManyPlayList                   `json:"playlist"`          // Not-Critical
	Thumbnail        *primitives.Image               `json:"thumbnail"`         // Not-Critical
	Banner           *primitives.Image               `json:"banner"`            // Not-Critical
	Tags             *primitives.TagList             `json:"tags"`              // Not-Critical
	SuggestedChannel *primitives.HashList            `json:"suggestedchannels"` // Not-Critical
	Content          []*CustomJSONMarshalContentList `json:"contentlist"`       // Not-Critical

	CreationTime time.Time `json:"creationtime"` // Not-Critical
}

func (a *CustomJSONMarshalChannel) IsSimilarTo(b CustomJSONMarshalChannel) bool {
	if a.RootChainID != b.RootChainID {
		return false
	}

	if a.ManagementChainID != b.ManagementChainID {
		return false
	}

	if a.ContentChainID != b.ContentChainID {
		return false
	}

	if a.LV1PublicKey != b.LV1PublicKey {
		return false
	}

	if a.LV2PublicKey != b.LV2PublicKey {
		return false
	}

	if a.LV3PublicKey != b.LV3PublicKey {
		return false
	}

	if !a.ChannelTitle.IsSameAs(b.ChannelTitle) {
		return false
	}

	if !a.Website.IsSameAs(b.Website) {
		return false
	}

	if !a.LongDescription.IsSameAs(b.LongDescription) {
		return false
	}

	if !a.ShortDescription.IsSameAs(b.ShortDescription) {
		return false
	}

	if !a.Playlist.IsSameAs(b.Playlist) {
		return false
	}

	if !a.Thumbnail.IsSameAs(b.Thumbnail) {
		return false
	}

	if !a.Tags.IsSameAs(b.Tags) {
		return false
	}

	if !a.SuggestedChannel.IsSameAs(b.SuggestedChannel) {
		return false
	}

	if len(a.Content) != len(b.Content) {
		return false
	}
	// TODO: Compare content lists

	return true
}

type CustomJSONMarshalContentList struct {
	ContentID string `json:"contentid"`
	Title     string `json:"title"`
}

func (a *Channel) ToCustomMarsalStruct() CustomJSONMarshalChannel {
	con := make([]*CustomJSONMarshalContentList, 0)
	for _, h := range a.Content.GetContents() {
		ci := new(CustomJSONMarshalContentList)
		ci.Title = h.ContentTitle.String()
		ci.ContentID = h.ContentID.String()
		con = append(con, ci)
		//hashList = append(hashList, h.ContentID.String())
		//titleList = append(titleList, h.ContentTitle.String())
	}

	custom := CustomJSONMarshalChannel{
		RootChainID:       a.RootChainID.String(),
		ManagementChainID: a.ManagementChainID.String(),
		ContentChainID:    a.ContentChainID.String(),
		LV1PublicKey:      a.LV1PublicKey.String(),
		LV2PublicKey:      a.LV2PublicKey.String(),
		LV3PublicKey:      a.LV3PublicKey.String(),
		ContentSingingKey: a.ContentSingingKey.String(),
		ChannelTitle:      &a.ChannelTitle,
		Website:           &a.Website,
		LongDescription:   &a.LongDescription,
		ShortDescription:  &a.ShortDescription,
		Playlist:          &a.Playlist,
		Thumbnail:         &a.Thumbnail,
		Banner:            &a.Banner,
		Tags:              &a.Tags,
		SuggestedChannel:  &a.SuggestedChannel,
		Content:           con,
		CreationTime:      a.CreationTime,
	}
	return custom
}

// CustomMarshalJSON reduces overhead of contentlist
func (a *Channel) CustomMarshalJSON() ([]byte, error) {
	custom := a.ToCustomMarsalStruct()
	return json.Marshal(custom)
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
	if (a == nil && b != nil) || (b == nil && a != nil) {
		return false
	}

	if !a.ChannelTitle.IsSameAs(&b.ChannelTitle) {
		//log.Println("Exit 0")
		return false
	}

	if !a.RootChainID.IsSameAs(&b.RootChainID) {
		//log.Println("Exit 1")
		return false
	}

	if !a.ManagementChainID.IsSameAs(&b.ManagementChainID) {
		//log.Println("Exit 1")
		return false
	}

	if !a.ContentChainID.IsSameAs(&b.ContentChainID) {
		//log.Println("Exit 2")
		return false
	}

	if !a.LV1PublicKey.IsSameAs(&b.LV1PublicKey) {
		//log.Println("Exit 3")
		return false
	}

	if !a.LV2PublicKey.IsSameAs(&b.LV2PublicKey) {
		//log.Println("Exit 4")
		return false
	}

	if !a.LV3PublicKey.IsSameAs(&b.LV3PublicKey) {
		//log.Println("Exit 5")
		return false
	}

	if !a.ContentSingingKey.IsSameAs(&b.ContentSingingKey) {
		//log.Println("Exit 6")
		return false
	}

	if !a.Website.IsSameAs(&b.Website) {
		//log.Printf("Exit 7. Web A: %s, Web B: %s\n", a.Website.String(), b.Website.String())
		return false
	}

	if !a.LongDescription.IsSameAs(&b.LongDescription) {
		//log.Println("Exit 8")
		return false
	}

	if !a.ShortDescription.IsSameAs(&b.ShortDescription) {
		//log.Println("Exit 9")
		return false
	}

	if !a.Playlist.IsSameAs(&b.Playlist) {
		//log.Println("Exit 10")
		return false
	}

	if !a.Thumbnail.IsSameAs(&b.Thumbnail) {
		//log.Println("Exit 11")
		return false
	}

	if !a.Banner.IsSameAs(&b.Banner) {
		//log.Println("Exit 12")
		return false
	}

	if !a.Tags.IsSameAs(&b.Tags) {
		//log.Println(a.Tags, b.Tags)
		//log.Println("Exit 13")
		return false
	}

	if !a.SuggestedChannel.IsSameAs(&b.SuggestedChannel) {
		//log.Println("Exit 14")
		return false
	}

	if !a.Content.IsSameAs(&b.Content) {
		//log.Println("Exit 15")
		//log.Println("DEBUG:", len(a.Content.GetContents()), len(b.Content.ContentList))
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
