package common

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type ManyPlayList struct {
	playlists []SinglePlayList
}

func (pl *ManyPlayList) GetPlaylists() []SinglePlayList {
	return pl.playlists
}

func (pl *ManyPlayList) Empty() bool {
	return len(pl.playlists) == 0
}

func RandomManyPlayList(max uint32) *ManyPlayList {
	p := new(ManyPlayList)
	u := random.RandomUInt32Between(0, max)

	p.playlists = make([]SinglePlayList, u)

	for i := range p.playlists {
		p.playlists[i] = *RandomSinglePlayList(max)
	}

	return p
}

func SmartRandomManyPlayList(max uint32, conts ContentList) *ManyPlayList {
	p := new(ManyPlayList)
	u := random.RandomUInt32Between(0, max)

	p.playlists = make([]SinglePlayList, u)

	for i := range p.playlists {
		p.playlists[i] = *SmartRandomSinglePlayList(max, conts)
	}

	return p
}

func (a *ManyPlayList) IsSameAs(b *ManyPlayList) bool {
	if len(a.playlists) != len(b.playlists) {
		return false
	}

	for i := range a.playlists {
		if !a.playlists[i].IsSameAs(&b.playlists[i]) {
			return false
		}
	}

	return true
}

func (p *ManyPlayList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data := primitives.Uint32ToBytes(uint32(len(p.playlists)))
	buf.Write(data)

	for i := range p.playlists {
		data, err := p.playlists[i].MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(data)
	}

	return buf.Next(buf.Len()), nil
}

func (p *ManyPlayList) UnmarshalBinary(data []byte) error {
	_, err := p.UnmarshalBinaryData(data)
	return err
}

func (p *ManyPlayList) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
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

	p.playlists = make([]SinglePlayList, u)
	var i uint32
	for i = 0; i < u; i++ {
		newData, err = p.playlists[i].UnmarshalBinaryData(newData)
		if err != nil {
			return data, err
		}
	}

	return
}

// SinglePlayList
type SinglePlayList struct {
	Title    primitives.Title
	Playlist primitives.HashList
}

func SmartRandomSinglePlayList(max uint32, contList ContentList) *SinglePlayList {
	p := new(SinglePlayList)

	conts := contList.GetContents()
	p.Playlist = *primitives.NewHashList()
	var i uint32
	for i = 0; i < max; i++ {
		if i >= uint32(len(conts)) {
			break
		}
		r := random.RandomIntBetween(0, 99)
		if r%2 == 0 {
			p.Playlist.AddHash(&conts[i].ContentID)
		}
	}
	p.Title = *primitives.RandomTitle()

	return p
}

func RandomSinglePlayList(max uint32) *SinglePlayList {
	p := new(SinglePlayList)

	p.Playlist = *primitives.RandomHashList(max)
	p.Title = *primitives.RandomTitle()

	return p
}

func (a *SinglePlayList) IsSameAs(b *SinglePlayList) bool {
	if !a.Playlist.IsSameAs(&b.Playlist) {
		return false
	}

	if !a.Title.IsSameAs(&b.Title) {
		return false
	}

	return true
}

func (p *SinglePlayList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data, err := p.Playlist.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = p.Title.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	return buf.Next(buf.Len()), nil
}

func (p *SinglePlayList) UnmarshalBinary(data []byte) error {
	_, err := p.UnmarshalBinaryData(data)
	return err
}

func (p *SinglePlayList) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data

	newData, err = p.Playlist.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = p.Title.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	return
}
