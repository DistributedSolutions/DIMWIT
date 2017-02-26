package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/DistributedSolutions/DIMWIT/common/primitives"
	"github.com/DistributedSolutions/DIMWIT/common/primitives/random"
)

var _ = log.Prefix()

type ManyPlayList struct {
	Playlists []SinglePlayList
}

func (pl *ManyPlayList) GetPlaylists() []SinglePlayList {
	return pl.Playlists
}

func (pl *ManyPlayList) Empty() bool {
	return len(pl.Playlists) == 0
}

func RandomManyPlayList(max uint32) *ManyPlayList {
	p := new(ManyPlayList)
	u := random.RandomUInt32Between(0, max)

	p.Playlists = make([]SinglePlayList, u)

	for i := range p.Playlists {
		p.Playlists[i] = *RandomSinglePlayList(max)
	}

	return p
}

func SmartRandomManyPlayList(max uint32, conts ContentList) *ManyPlayList {
	p := new(ManyPlayList)
	u := random.RandomUInt32Between(0, max)

	p.Playlists = make([]SinglePlayList, u)

	for i := range p.Playlists {
		p.Playlists[i] = *SmartRandomSinglePlayList(max, conts)
	}

	return p
}

func (a *ManyPlayList) Combine(b *ManyPlayList) *ManyPlayList {
	pl := append(a.Playlists, b.Playlists...)
	x := new(ManyPlayList)
	x.Playlists = pl
	return x
}

func (a *ManyPlayList) IsSameAs(b *ManyPlayList) bool {
	if len(a.Playlists) != len(b.Playlists) {
		// log.Println("[Playlist] Exit 1", len(a.Playlists), len(b.Playlists))
		return false
	}

	for i := range a.Playlists {
		if !a.Playlists[i].IsSameAs(&b.Playlists[i]) {
			// log.Println("[Playlist] Exit 2")
			return false
		}
	}

	return true
}

func (p *ManyPlayList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data := primitives.Uint32ToBytes(uint32(len(p.Playlists)))
	buf.Write(data)

	for i := range p.Playlists {
		data, err := p.Playlists[i].MarshalBinary()
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

	p.Playlists = make([]SinglePlayList, u)
	var i uint32
	for i = 0; i < u; i++ {
		newData, err = p.Playlists[i].UnmarshalBinaryData(newData)
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

func (p *SinglePlayList) MarshalJSON() ([]byte, error) {
	list := make([]string, 0)
	for _, h := range p.Playlist.GetHashes() {
		list = append(list, h.String())
	}
	return json.Marshal(&struct {
		Title    string   `json:"title"`
		Playlist []string `json:"playlist"`
	}{
		Title:    p.Title.String(),
		Playlist: list,
	})
}

func (p *SinglePlayList) UnmarshalJSON(b []byte) error {
	obj := new(struct {
		Title    string   `json:"title"`
		Playlist []string `json:"playlist"`
	})

	if err := json.Unmarshal(b, obj); err != nil {
		return err
	}

	p.Title.SetString(obj.Title)
	for _, h := range obj.Playlist {
		hash, err := primitives.HexToHash(h)
		if err != nil {
			return err
		}
		p.Playlist.AddHash(hash)
	}
	return nil
}
