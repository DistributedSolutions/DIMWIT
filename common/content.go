package common

import (
	"bytes"
	"fmt"

	"github.com/DistributedSolutions/DIMWIT/common/primitives"
)

type ContentList struct {
	Length      int
	ContentList []Content
}

type Content struct {
	// Content Metadada
	Type             byte
	RootChainID      primitives.Hash
	InfoHash         primitives.InfoHash
	LongDescription  primitives.LongDescription
	ShortDescription primitives.ShortDescription
	ActionFiles      primitives.FileList
	Thumbnail        primitives.Image
	Series           byte
	// There can be lots of parts
	Part [2]byte
	Tags primitives.TagList

	// Torrent Metadata
	Tracker  primitives.TrackerList
	FileList primitives.FileList
}

func NewContent() *Content {
	c := new(Content)
	//c.RootChainID = primitives.New
	return c
}

func (c *Content) MarshalBinary() (data []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("A panic has occurred while marshaling: %s", r)
			return
		}
	}()

	//We won't know exactly where it fails because we don't catch every error, but this looks so much cleaner,
	// and it if fails, then returning where it fails only helps with debugging. It should only fail if
	// you made a *Content on your own and not through the proper functions

	buf := new(bytes.Buffer)

	buf.Write([]byte{c.Type})

	data, err = c.RootChainID.MarshalBinary()
	buf.Write(data)

	data, err = c.InfoHash.MarshalBinary()
	buf.Write(data)

	data, err = c.LongDescription.MarshalBinary()
	buf.Write(data)

	data, err = c.ShortDescription.MarshalBinary()
	buf.Write(data)

	data, err = c.ActionFiles.MarshalBinary()
	buf.Write(data)

	data, err = c.Thumbnail.MarshalBinary()
	buf.Write(data)

	buf.Write([]byte{c.Series})

	buf.Write(c.Part[:])

	data, err = c.Tags.MarshalBinary()
	buf.Write(data)

	data, err = c.Tracker.MarshalBinary()
	buf.Write(data)

	data, err = c.FileList.MarshalBinary()
	buf.Write(data)

	if err != nil {
		return nil, fmt.Errorf("Failed to marshal content")
	}

	return buf.Next(buf.Len()), err
}

func (c *Content) UnmarshalBinaryData(data []byte) (newData []byte, err error) {

	return nil, nil
}
