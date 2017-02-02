package common

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

type ManyPlayList struct {
	length    uint32
	title     primitives.Title
	playlists []primitives.HashList
}

func (pl *ManyPlayList) GetPlaylists() []primitives.HashList {
	return pl.playlists
}

func RandomManyPlayList(max uint32) *ManyPlayList {
	p := new(ManyPlayList)
	u := random.RandomUInt32Between(0, max)

	p.length = u
	p.playlists = make([]primitives.HashList, u)

	for i := range p.playlists {
		p.playlists[i] = *primitives.RandomHashList(max)
	}

	p.title = *primitives.RandomTitle()

	return p
}

func (a *ManyPlayList) IsSameAs(b *ManyPlayList) bool {
	if a.length != b.length {
		return false
	}

	for i := range a.playlists {
		if !a.playlists[i].IsSameAs(&b.playlists[i]) {
			return false
		}
	}

	if !a.title.IsSameAs(&b.title) {
		return false
	}

	return true
}

func (p *ManyPlayList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data, err := primitives.Uint32ToBytes(p.length)
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	for i := range p.playlists {
		data, err = p.playlists[i].MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(data)
	}

	data, err = p.title.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

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
	p.length = u
	newData = newData[4:]

	p.playlists = make([]primitives.HashList, u)
	var i uint32
	for i = 0; i < u; i++ {
		newData, err = p.playlists[i].UnmarshalBinaryData(newData)
		if err != nil {
			return data, err
		}
	}

	newData, err = p.title.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	return
}
