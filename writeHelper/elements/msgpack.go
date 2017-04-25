package elements

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/DistributedSolutions/DIMWIT/common"
	"github.com/DistributedSolutions/DIMWIT/common/primitives"
)

//go:generate msgp
type ManageChainMetaDataBytes struct {
	Website           []byte `msg:"website"`
	LongDescription   []byte `msg:"longdesc"`
	ShortDescription  []byte `msg:"shortdesc"`
	Playlist          []byte `msg:"playlist"`
	Thumbnail         []byte `msg:"thumbnail"`
	Banner            []byte `msg:"banner"`
	ChannelTags       []byte `msg:"chantags"`
	SuggestedChannels []byte `msg:"sugchans"`
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
